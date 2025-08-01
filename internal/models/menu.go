package models

import (
	"encoding/json"
	"strings"
	"time"

	"gin-admin/internal/configs"
)

const (
	MenuStatus_DISABLED = "disabled"
	MenuStatus_ENABLED  = "enabled"

	MenuType_CATALOG = "catalog"
	MenuType_MENU    = "menu"
	MenuType_BUTTON  = "button"
)

// Menu management for SYS
type Menu struct {
	ID         string         `json:"id" gorm:"size:20;primarykey;"`                   // Unique ID
	Name       string         `json:"name" gorm:"size:128;index"`                      // Display name of menu
	Type       string         `json:"type" gorm:"size:20;index"`                       // Type of menu (catalog, menu, button)
	Method     string         `json:"method" gorm:"size:20;index;"`                    // Http method of resource
	Path       string         `json:"path" gorm:"size:255;"`                           // Access path of menu
	Component  string         `json:"component" gorm:"size:255;"`                      // Component path of view
	Status     string         `json:"status" gorm:"size:20;index"`                     // Status of menu (enabled, disabled)
	Redirect   string         `json:"redirect" gorm:"size:255;not null;default:''"`    // Redirect path of menu
	ParentID   string         `json:"parentId" gorm:"size:20;index;"`                  // Parent ID (From Menu.ID)
	ParentPath string         `json:"parentPath" gorm:"size:255;index;"`               // Parent path (split by .)
	Rank       int            `json:"rank" gorm:"column:rank;index;"`                  // Rank for sorting (Order by desc)
	Title      string         `json:"title" gorm:"size:1024"`                          // Menu title
	Extra      map[string]any `json:"-" gorm:"type:text;serializer:json;default:'{}'"` // Meta data
	CreatedAt  time.Time      `json:"createdAt" gorm:"index;"`                         // Create time
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"index;"`                         // Update time

	Children *Menus `json:"children" gorm:"-"` // Child menus
	Roles    Roles  `json:"roles" gorm:"many2many:role_menus;"`
}

func (a Menu) TableName() string {
	return configs.C.FormatTableName("menu")
}

// 自定义JSON反序列化
func (a *Menu) UnmarshalJSON(data []byte) error {
	// 使用别名类型避免递归调用
	type Alias Menu
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	// 解析固定字段
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	fields := make(map[string]struct{})
	fields["id"] = struct{}{}
	fields["name"] = struct{}{}
	fields["type"] = struct{}{}
	fields["method"] = struct{}{}
	fields["path"] = struct{}{}
	fields["component"] = struct{}{}
	fields["status"] = struct{}{}
	fields["redirect"] = struct{}{}
	fields["parentId"] = struct{}{}
	fields["parentPath"] = struct{}{}
	fields["rank"] = struct{}{}
	fields["title"] = struct{}{}
	fields["createdAt"] = struct{}{}
	fields["updatedAt"] = struct{}{}
	fields["children"] = struct{}{}
	fields["roles"] = struct{}{}

	// 解析整个JSON到map
	var rawMap map[string]any
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// 初始化Extra map
	a.Extra = make(map[string]any)

	// 移除固定字段，保留动态字段
	for key := range rawMap {
		if _, exist := fields[key]; !exist {
			a.Extra[key] = rawMap[key]
		}
	}

	return nil
}

// 自定义JSON序列化
func (a Menu) MarshalJSON() ([]byte, error) {

	raw := make(map[string]any)
	for key, value := range a.Extra {
		raw[key] = value
	}

	raw["id"] = a.ID
	raw["name"] = a.Name
	raw["type"] = a.Type
	raw["method"] = a.Method
	raw["path"] = a.Path
	raw["component"] = a.Component
	raw["status"] = a.Status
	raw["redirect"] = a.Redirect
	raw["parentId"] = a.ParentID
	raw["parentPath"] = a.ParentPath
	raw["rank"] = a.Rank
	raw["title"] = a.Title
	raw["createdAt"] = a.CreatedAt
	raw["updatedAt"] = a.UpdatedAt
	raw["children"] = a.Children
	raw["roles"] = a.Roles

	return json.Marshal(raw)
}

// Defining the slice of `Menu` struct.
type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	if a[i].Rank == a[j].Rank {
		return a[i].CreatedAt.Unix() > a[j].CreatedAt.Unix()
	}
	return a[i].Rank > a[j].Rank
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Menus) ToIDMapper() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// collect all parent IDs of menu list
func (a Menus) ParentIDs() []string {
	parentIDs := []string{}
	cacher := map[string]struct{}{}
	for _, item := range a {
		if _, ok := cacher[item.ID]; ok {
			continue
		}
		cacher[item.ID] = struct{}{}
		if pp := item.ParentPath; pp != "" {
			for _, pid := range strings.Split(pp, ".") {
				if pid == "" {
					continue
				}
				if _, ok := cacher[pid]; ok {
					continue
				}
				parentIDs = append(parentIDs, pid)
				cacher[pid] = struct{}{}
			}
		}
	}
	return parentIDs
}

func (a Menus) ToTree() Menus {
	var list Menus
	m := a.ToIDMapper()
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if parent, ok := m[item.ParentID]; ok {
			if parent.Children == nil {
				children := Menus{item}
				parent.Children = &children
			} else {
				*parent.Children = append(*parent.Children, item)
			}
		}
	}
	return list
}

func (m Menus) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("[]"), nil
	}

	type M Menus
	copy := M(m)

	return json.Marshal(copy)
}
