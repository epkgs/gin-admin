package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gin-admin/internal/configs"
	"gin-admin/internal/defines"
	"gin-admin/internal/dtos"
	"gin-admin/internal/errorx"
	"gin-admin/internal/models"
	"gin-admin/internal/repositories"
	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/encoding/json"
	"gin-admin/pkg/encoding/yaml"
	"gin-admin/pkg/gormx"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/randx"

	"github.com/epkgs/object"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const treePathDelimiter = "."

// Menu management for SYS
type Menu struct {
	Cacher       cachex.Cacher
	MenuRepo     *repositories.Menu
	MenuRoleRepo *repositories.MenuRole
	UserRoleRepo *repositories.UserRole
}

func NewMenu(app types.AppContext) *Menu {
	return &Menu{
		Cacher:       app.Cacher(),
		MenuRepo:     repositories.NewMenu(app.DB()),
		MenuRoleRepo: repositories.NewMenuRole(app.DB()),
		UserRoleRepo: repositories.NewUserRole(app.DB()),
	}
}

func (a *Menu) InitIfNeed(ctx context.Context) error {
	if configs.C.Menu.File == "" {
		return nil
	}

	count, err := a.MenuRepo.Count(ctx)
	if err != nil {
		return errorx.WrapGormError(ctx, err)
	}

	if count > 0 {
		logger.Info(ctx, "Menu database is not empty, skip init menu data.")
		return nil // 已有数据就跳过
	}

	if err := a.initFromFile(ctx, configs.C.Menu.File); err != nil {
		logger.Error(ctx, "failed to init menu data", err, map[string]any{"file": configs.C.Menu.File})
	}

	return a.syncToCasbin(ctx)
}

func (a *Menu) initFromFile(ctx context.Context, menuFile string) error {

	menus := models.Menus{}

	var tmpMenus models.Menus

	f, err := os.ReadFile(menuFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Warn(ctx, "Menu data file not found, skip init menu data from file", map[string]any{"file": menuFile})
			return nil
		}
		return err
	}

	if ext := filepath.Ext(menuFile); ext == ".json" {
		if err := json.Unmarshal(f, &tmpMenus); err != nil {
			return errors.Wrapf(err, "Unmarshal JSON file '%s' failed", menuFile)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(f, &tmpMenus); err != nil {
			return errors.Wrapf(err, "Unmarshal YAML file '%s' failed", menuFile)
		}
	} else {
		return errors.Errorf("Unsupported file type '%s'", ext)
	}

	menus = append(menus, tmpMenus...)

	return a.upsert(ctx, menus, nil)
}

func (a *Menu) upsert(ctx context.Context, items models.Menus, parent *models.Menu) error {
	total := len(items)

	for i, item := range items {
		var parentID string
		if parent != nil {
			parentID = parent.ID
		}

		var (
			menuItem *models.Menu
			err      error
		)

		if item.ID != "" {
			menuItem, err = a.MenuRepo.Get(ctx, item.ID)
		} else if item.Name != "" {
			menuItem, err = a.MenuRepo.GetChildByName(ctx, parentID, item.Name)
		}

		if err != nil {
			return errorx.WrapGormError(ctx, err)
		}

		if item.Status == "" {
			item.Status = models.MenuStatus_ENABLED
		}

		if menuItem != nil {
			var md object.Metadata
			object.Assign(menuItem, item, func(c *object.AssignConfig) {
				c.Metadata = &md
			})

			if len(md.Keys) > 0 { // changed
				if err := a.MenuRepo.Update(ctx, menuItem, gormx.WithSelect(md.Keys), gormx.WithOmit("Children")); err != nil {
					return err
				}
			}

		} else {
			if item.ID == "" {
				item.ID = randx.NewXID()
			}
			if item.Rank == 0 {
				item.Rank = total - i
			}
			item.ParentID = parentID
			if parent != nil {
				item.ParentPath = parent.ParentPath + parentID + treePathDelimiter
			}
			menuItem = item
			if err := a.MenuRepo.Create(ctx, item, gormx.WithOmit("Children")); err != nil {
				return err
			}
		}

		if item.Children != nil {
			if err := a.upsert(ctx, *item.Children, menuItem); err != nil {
				return err
			}
		}
	}
	return nil
}

// List menus from the data access object based on the provided parameters and options.
func (a *Menu) List(ctx context.Context, req dtos.MenuListReq) (*dtos.List[*models.Menu], error) {

	option := func(db *gorm.DB) *gorm.DB {
		if v := req.InIDs; len(v) > 0 {
			db = db.Where("id IN ?", v)
		}
		if v := req.LikeName; len(v) > 0 {
			db = db.Where("name LIKE ?", "%"+v+"%")
		}
		if v := req.Status; len(v) > 0 {
			db = db.Where("status = ?", v)
		}
		if v := req.ParentID; len(v) > 0 {
			db = db.Where("parent_id = ?", v)
		}
		if v := req.ParentPathPrefix; len(v) > 0 {
			db = db.Where("parent_path LIKE ?", v+"%")
		}
		if v := req.UserID; len(v) > 0 {
			userRoleQuery := a.UserRoleRepo.DB().Model(new(models.UserRole)).Where("user_id = ?", v).Select("role_id")
			menuRoleQuery := a.MenuRoleRepo.DB().Model(new(models.MenuRole)).Where("role_id IN (?)", userRoleQuery).Select("menu_id")
			db = db.Where("id IN (?)", menuRoleQuery)
		}
		if v := req.RoleID; len(v) > 0 {
			menuRoleQuery := a.MenuRoleRepo.DB().Model(new(models.MenuRole)).Where("role_id = ?", v).Select("menu_id")
			db = db.Where("id IN (?)", menuRoleQuery)
		}
		if v := req.Type; len(v) > 0 {
			db = db.Where("type = ?", v)
		} else {
			if v := req.WithResources; !v {
				db = db.Where("type != ?", models.MenuType_BUTTON)
			}
		}

		return db
	}

	list, err := a.MenuRepo.Find(ctx, option, gormx.WithPage(req.Page, req.Limit), gormx.WithOrder("rank", "desc"), gormx.WithOrder("created_at", "desc"))
	if err != nil {
		return nil, errorx.WrapGormError(ctx, err)
	}

	count, err := a.MenuRepo.Count(ctx, option)
	if err != nil {
		return nil, errorx.WrapGormError(ctx, err)
	}

	if req.LikeName != "" {
		list, err = a.appendChildren(ctx, list)
		if err != nil {
			return nil, errorx.WrapGormError(ctx, err)
		}
	}

	return dtos.NewList(list, req.Page, req.Limit, count), nil
}

func (a *Menu) appendChildren(ctx context.Context, data models.Menus) (models.Menus, error) {
	if len(data) == 0 {
		return data, nil
	}

	dataCache := map[string]struct{}{}
	// init cache
	for _, menu := range data {
		dataCache[menu.ID] = struct{}{}
	}

	appendData := func(child *models.Menu) {
		if _, exist := dataCache[child.ID]; exist {
			return
		}

		dataCache[child.ID] = struct{}{}
		data = append(data, child)
	}

	for _, item := range data {
		children, err := a.MenuRepo.Find(ctx, gormx.WithWhere("parent_path LIKE ?", item.ParentPath+item.ID+treePathDelimiter+"%"))
		if err != nil {
			return nil, errorx.WrapGormError(ctx, err)
		}
		for _, child := range children {
			appendData(child)
		}
	}

	if parentIDs := data.ParentIDs(); len(parentIDs) > 0 {
		parents, err := a.MenuRepo.Find(ctx, gormx.WithWhere("id IN (?)", parentIDs))
		if err != nil {
			return nil, errorx.WrapGormError(ctx, err)
		}
		for _, p := range parents {
			appendData(p)
		}
	}
	sort.Sort(data)

	return data, nil
}

// Get the specified menu from the data access object.
func (a *Menu) Get(ctx context.Context, id string) (*models.Menu, error) {
	menu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return nil, errorx.WrapGormError(ctx, err)
	}

	children, err := a.MenuRepo.Find(ctx, gormx.WithWhere("parent_id = ?", menu.ID))
	if err != nil {
		return nil, errorx.WrapGormError(ctx, err)
	}

	var menus models.Menus = children
	menu.Children = &menus

	return menu, nil
}

// Create a new menu in the data access object.
func (a *Menu) Create(ctx context.Context, req *dtos.MenuCreateReq) (*models.Menu, error) {
	if configs.C.Menu.DenyOperate {
		return nil, errorx.ErrBadRequest.New(ctx)
	}

	menu := &models.Menu{
		ID:        randx.NewXID(),
		CreatedAt: time.Now(),
	}

	if parentID := req.ParentID; parentID != "" {
		parent, err := a.MenuRepo.Get(ctx, parentID)
		if err != nil {
			return nil, errorx.WrapGormError(ctx, err)
		}

		menu.ParentPath = parent.ParentPath + parent.ID + treePathDelimiter
	}

	if err := object.Assign(menu, req); err != nil {
		return nil, err
	}

	if err := a.MenuRepo.Create(ctx, menu); err != nil {
		return nil, errorx.WrapGormError(ctx, err)
	}

	return menu, nil
}

// Update the specified menu in the data access object.
func (a *Menu) Update(ctx context.Context, id string, req *dtos.MenuUpdateReq) error {
	if configs.C.Menu.DenyOperate {
		return errorx.ErrBadRequest.New(ctx)
	}

	menu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return errorx.WrapGormError(ctx, err)
	}

	oldParentPath := menu.ParentPath
	oldStatus := menu.Status
	var childData models.Menus
	if req.ParentID != nil && menu.ParentID != *req.ParentID {
		if parentID := *req.ParentID; parentID != "" {
			parent, err := a.MenuRepo.Get(ctx, parentID)
			if err != nil {
				return errorx.WrapGormError(ctx, err)
			}
			menu.ParentPath = parent.ParentPath + parent.ID + treePathDelimiter
		} else {
			menu.ParentPath = ""
		}

		res, err := a.MenuRepo.Find(ctx, gormx.WithWhere("parent_path LIKE ?", oldParentPath+menu.ID+treePathDelimiter+"%"), gormx.WithSelect("id", "parent_path"))
		if err != nil {
			return errorx.WrapGormError(ctx, err)
		}
		childData = res
	}

	if err := object.Assign(menu, req); err != nil {
		return errorx.ErrInternal.New(ctx).Wrap(err)
	}

	err = a.MenuRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if req.Status != nil && oldStatus != *req.Status {
			oldPath := oldParentPath + menu.ID + treePathDelimiter
			if err := a.MenuRepo.UpdateStatusByParentPath(ctx, oldPath, *req.Status); err != nil {
				return err
			}
		}

		for _, child := range childData {
			oldPath := oldParentPath + menu.ID + treePathDelimiter
			newPath := menu.ParentPath + menu.ID + treePathDelimiter
			err := a.MenuRepo.UpdateParentPath(ctx, child.ID, strings.Replace(child.ParentPath, oldPath, newPath, 1))
			if err != nil {
				return err
			}
		}

		if err := a.MenuRepo.Update(ctx, menu); err != nil {
			return err
		}

		if menu.Type != models.MenuType_MENU {
			if err := a.MenuRepo.DeleteChildrenOfButton(ctx, menu.ParentID); err != nil {
				return err
			}
		}
		return a.syncToCasbin(ctx)
	})

	return errorx.WrapGormError(ctx, err)
}

// Delete the specified menu from the data access object.
func (a *Menu) Delete(ctx context.Context, id string) error {
	if configs.C.Menu.DenyOperate {
		return errorx.ErrBadRequest.New(ctx)
	}

	menu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return errorx.WrapGormError(ctx, err)
	}

	children, err := a.MenuRepo.Find(ctx, gormx.WithWhere("parent_path LIKE ?", menu.ParentPath+menu.ID+treePathDelimiter+"%"), gormx.WithSelect("id"))
	if err != nil {
		return errorx.WrapGormError(ctx, err)
	}

	err = a.MenuRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := a.delete(ctx, id); err != nil {
			return err
		}

		for _, child := range children {
			if err := a.delete(ctx, child.ID); err != nil {
				return err
			}
		}

		return a.syncToCasbin(ctx)
	})

	return errorx.WrapGormError(ctx, err)
}

func (a *Menu) delete(ctx context.Context, id string) error {
	if err := a.MenuRepo.Delete(ctx, id); err != nil {
		return err
	}
	if err := a.MenuRoleRepo.DeleteByMenuID(ctx, id); err != nil {
		return err
	}
	return nil
}

func (a *Menu) syncToCasbin(ctx context.Context) error {
	return a.Cacher.Set(ctx, defines.CacheNSForRole, defines.CacheKeyForSyncToCasbin, fmt.Sprintf("%d", time.Now().Unix()))
}
