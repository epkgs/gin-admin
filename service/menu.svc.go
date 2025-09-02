package service

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gin-admin/errorx"
	"gin-admin/locales"
	"gin-admin/model/bo"
	"gin-admin/model/dto"
	"gin-admin/model/po"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/encoding/json"
	"gin-admin/pkg/encoding/yaml"
	"gin-admin/pkg/logger"
	"gin-admin/pkg/randx"
	"gin-admin/types"

	"github.com/epkgs/i18n/errors"

	"github.com/epkgs/object"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

const (
	gTreePathDelimiter = "."
)

// Menu management for SYS
type Menu struct {
	app     types.AppContext
	cacher  cachex.Cacher
	roleSvc *Role
	q       *bo.Query
}

func NewMenu(app types.AppContext) *Menu {
	return &Menu{
		app:     app,
		cacher:  app.Cacher(),
		roleSvc: NewRole(app),
		q:       bo.Use(app.DB()),
	}
}

func (a *Menu) InitIfNeed(ctx context.Context) error {
	if a.app.Config().Menu.File == "" {
		return nil
	}

	m := a.q.Menu

	count, err := m.WithContext(ctx).Count()
	if err != nil {
		return errorx.WrapGormError(err)
	}

	if count > 0 {
		logger.Info(ctx, "Menu database is not empty, skip init menu data.")
		return nil // 已有数据就跳过
	}

	if err := a.initFromFile(ctx, a.app.Config().Menu.File); err != nil {
		logger.Error(ctx, "failed to init menu data", err, map[string]any{"file": a.app.Config().Menu.File})
	}

	return a.roleSvc.RefreshUpdateTime(ctx)
}

func (a *Menu) initFromFile(ctx context.Context, menuFile string) error {

	menus := po.Menus{}

	var tmpMenus po.Menus

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
			return errorx.ErrInternalServerError.WithMsg(locales.Def.Str("Unmarshal JSON file '%s' failed", menuFile))
		}
	} else if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(f, &tmpMenus); err != nil {
			return errorx.ErrInternalServerError.WithMsg(locales.Def.Str("unmarshal YAML file '%s' failed", menuFile))
		}
	} else {
		return errorx.ErrBadRequest.WithMsg(locales.Def.Str("unsupported file type '%s'", ext))
	}

	menus = append(menus, tmpMenus...)

	return a.upsert(ctx, menus, nil)
}

func (a *Menu) upsert(ctx context.Context, items po.Menus, parent *po.Menu) error {
	total := len(items)

	for i, item := range items {
		var parentID string
		if parent != nil {
			parentID = parent.ID
		}

		var (
			menu *po.Menu
			err  error
		)

		m := a.q.Menu

		if item.ID != "" {
			menu, err = m.WithContext(ctx).Where(m.ID.Eq(item.ID)).First()
		} else if item.Name != "" {
			menu, err = m.WithContext(ctx).Where(m.Name.Eq(item.Name), m.ParentID.Eq(parentID)).First()
		}

		if err != nil {
			return errorx.WrapGormError(err)
		}

		if item.Status == "" {
			item.Status = po.MenuStatus_ENABLED
		}

		if menu != nil {
			var md object.Metadata
			object.Assign(menu, item, func(c *object.AssignConfig) {
				c.Metadata = &md
			})

			if len(md.Keys) > 0 { // changed
				if _, err := m.WithContext(ctx).Updates(menu); err != nil {
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
				item.ParentPath = parent.ParentPath + parentID + gTreePathDelimiter
			}
			menu = item

			if err := m.WithContext(ctx).Create(menu); err != nil {
				return err
			}
		}

		if item.Children != nil {
			if err := a.upsert(ctx, item.Children, menu); err != nil {
				return err
			}
		}
	}
	return nil
}

// List menus from the data access object based on the provided parameters and options.
func (a *Menu) List(ctx context.Context, req dto.MenuListReq) (*dto.List[*po.Menu], error) {

	m := a.q.Menu

	scope := func(d gen.Dao) gen.Dao {

		if v := req.InIDs; len(v) > 0 {
			d = d.Where(m.ID.In(v...))
		}
		if v := req.LikeName; len(v) > 0 {
			d = d.Where(m.Name.Like("%" + v + "%"))
		}
		if v := req.Status; len(v) > 0 {
			d = d.Where(m.Status.Eq(v))
		}
		if v := req.ParentID; len(v) > 0 {
			d = d.Where(m.ParentID.Eq(v))
		}
		if v := req.ParentPathPrefix; len(v) > 0 {
			d = d.Where(m.ParentPath.Like(v + "%"))
		}
		if v := req.UserID; len(v) > 0 {

			db := a.app.DB()

			subquery := db.Raw(
				"SELECT menu_id FROM r_role_menus WHERE role_id IN (?)",
				db.Raw("SELECT role_id FROM r_role_users WHERE user_id = ?", v),
			)

			expr := field.ContainsSubQuery([]field.Expr{m.ID}, subquery)

			d = d.Where(expr)
		}
		if v := req.RoleID; len(v) > 0 {
			subquery := a.app.DB().Raw(
				"SELECT menu_id FROM r_role_menus WHERE role_id = ? ",
				v,
			)
			expr := field.ContainsSubQuery([]field.Expr{m.ID}, subquery)
			d = d.Where(expr)
		}
		if v := req.Type; len(v) > 0 {
			d = d.Where(m.Type.Eq(v))
		} else {
			if v := req.WithResources; !v {
				d = d.Where(m.Type.Neq(po.MenuType_BUTTON))
			}
		}

		return d
	}

	list, err := m.WithContext(ctx).Scopes(scope, req.PageScope()).Order(m.Rank.Desc()).Order(m.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	count, err := m.WithContext(ctx).Scopes(scope).Count()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	if req.LikeName != "" {
		list, err = a.appendChildren(ctx, list)
		if err != nil {
			return nil, errorx.WrapGormError(err)
		}
	}

	return dto.NewList(list, &dto.Pager{
		Page:  req.Page,
		Limit: req.Limit,
		Total: count,
	}), nil
}

func (a *Menu) appendChildren(ctx context.Context, data po.Menus) (po.Menus, error) {
	if len(data) == 0 {
		return data, nil
	}

	dataCache := map[string]struct{}{}
	// init cache
	for _, menu := range data {
		dataCache[menu.ID] = struct{}{}
	}

	appendData := func(child *po.Menu) {
		if _, exist := dataCache[child.ID]; exist {
			return
		}

		dataCache[child.ID] = struct{}{}
		data = append(data, child)
	}

	m := a.q.Menu

	for _, item := range data {
		children, err := m.WithContext(ctx).Where(m.ParentPath.Like(item.ParentPath + item.ID + gTreePathDelimiter + "%")).Find()
		if err != nil {
			return nil, errorx.WrapGormError(err)
		}
		for _, child := range children {
			appendData(child)
		}
	}

	if parentIDs := data.ParentIDs(); len(parentIDs) > 0 {
		parents, err := m.WithContext(ctx).Where(m.ID.In(parentIDs...)).Find()
		if err != nil {
			return nil, errorx.WrapGormError(err)
		}
		for _, p := range parents {
			appendData(p)
		}
	}
	sort.Sort(data)

	return data, nil
}

// Get the specified menu from the data access object.
func (a *Menu) Get(ctx context.Context, id string) (*po.Menu, error) {
	m := a.q.Menu

	menu, err := m.WithContext(ctx).Get(id)
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	children, err := m.WithContext(ctx).Where(m.ParentID.Eq(menu.ID)).Find()
	if err != nil {
		return nil, errorx.WrapGormError(err)
	}

	menu.Children = children

	return menu, nil
}

// Create a new menu in the data access object.
func (a *Menu) Create(ctx context.Context, req *dto.MenuCreateReq) (*po.Menu, error) {
	if a.app.Config().Menu.DenyOperate {
		return nil, errorx.ErrBadRequest
	}

	menu := &po.Menu{
		ID:        randx.NewXID(),
		CreatedAt: time.Now(),
	}

	m := a.q.Menu

	if parentID := req.ParentID; parentID != "" {
		parent, err := m.WithContext(ctx).Get(parentID)
		if err != nil {
			return nil, errorx.WrapGormError(err)
		}

		menu.ParentPath = parent.ParentPath + parent.ID + gTreePathDelimiter
	}

	if err := object.Assign(menu, req, func(c *object.AssignConfig) {
		c.IncludeIgnoreFields = true
	}); err != nil {
		return nil, err
	}

	if err := m.WithContext(ctx).Create(menu); err != nil {
		return nil, errorx.WrapGormError(err)
	}

	return menu, nil
}

// Update the specified menu in the data access object.
func (a *Menu) Update(ctx context.Context, id string, req *dto.MenuUpdateReq) error {
	if a.app.Config().Menu.DenyOperate {
		return errorx.ErrBadRequest
	}

	m := a.q.Menu

	menu, err := m.WithContext(ctx).Get(id)
	if err != nil {
		return errorx.WrapGormError(err)
	}

	oldParentPath := menu.ParentPath
	oldStatus := menu.Status
	var childData po.Menus
	if req.ParentID != nil && menu.ParentID != *req.ParentID {
		if parentID := *req.ParentID; parentID != "" {
			parent, err := m.WithContext(ctx).Get(parentID)
			if err != nil {
				return errorx.WrapGormError(err)
			}
			menu.ParentPath = parent.ParentPath + parent.ID + gTreePathDelimiter
		} else {
			menu.ParentPath = ""
		}

		res, err := m.WithContext(ctx).
			Where(m.ParentPath.Like(oldParentPath+menu.ID+gTreePathDelimiter+"%")).
			Select(m.ID, m.ParentPath).
			Find()
		if err != nil {
			return errorx.WrapGormError(err)
		}
		childData = res
	}

	if err := object.Assign(menu, req); err != nil {
		return errorx.ErrInternalServerError.Wrap(err)
	}

	err = a.q.Transaction(func(tx *bo.Query) error {
		if req.Status != nil && oldStatus != *req.Status {
			oldPath := oldParentPath + menu.ID + gTreePathDelimiter
			_, err := m.WithContext(ctx).
				Where(m.ParentPath.Like(oldPath+"%")).
				Update(m.Status, *req.Status)
			if err != nil {
				return err
			}
		}

		for _, child := range childData {
			oldPath := oldParentPath + menu.ID + gTreePathDelimiter
			newPath := menu.ParentPath + menu.ID + gTreePathDelimiter
			_, err := m.WithContext(ctx).
				Where(m.ID.Eq(child.ID)).
				Update(m.ParentPath, strings.Replace(child.ParentPath, oldPath, newPath, 1))
			if err != nil {
				return err
			}
		}

		if _, err := m.WithContext(ctx).Updates(menu); err != nil {
			return err
		}

		if menu.Type != po.MenuType_MENU {
			_, err := m.WithContext(ctx).
				Where(m.ParentID.Eq(menu.ParentID)).
				Where(m.Type.Eq(po.MenuType_BUTTON)).
				Delete()
			if err != nil {
				return err
			}
		}
		return a.roleSvc.RefreshUpdateTime(ctx)
	})

	return errorx.WrapGormError(err)
}

// Delete the specified menu from the data access object.
func (a *Menu) Delete(ctx context.Context, id string) error {
	if a.app.Config().Menu.DenyOperate {
		return errorx.ErrBadRequest
	}

	m := a.q.Menu

	menu, err := m.WithContext(ctx).Get(id)
	if err != nil {
		return errorx.WrapGormError(err)
	}

	children, err := m.WithContext(ctx).Where(m.ParentPath.Like(menu.ParentPath + menu.ID + gTreePathDelimiter + "%")).
		Select(m.ID).
		Find()
	if err != nil {
		return errorx.WrapGormError(err)
	}

	err = a.q.Transaction(func(tx *bo.Query) error {
		if err := a.delete(ctx, id); err != nil {
			return err
		}

		for _, child := range children {
			if err := a.delete(ctx, child.ID); err != nil {
				return err
			}
		}

		return a.roleSvc.RefreshUpdateTime(ctx)
	})

	return errorx.WrapGormError(err)
}

func (a *Menu) delete(ctx context.Context, id string) error {
	m := a.q.Menu

	if err := m.Roles.Model(&po.Menu{ID: id}).Clear(); err != nil {
		return err
	}

	// err = a.app.DB().Exec("DELETE FROM r_role_menus WHERE menu_id = ?", id).Error
	// if err != nil {
	// 	return err
	// }

	_, err := m.WithContext(ctx).Where(m.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}

	return nil
}
