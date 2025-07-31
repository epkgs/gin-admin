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
	ID         string    `json:"id" gorm:"size:20;primarykey;"`                // Unique ID
	Name       string    `json:"name" gorm:"size:128;index"`                   // Display name of menu
	Type       string    `json:"type" gorm:"size:20;index"`                    // Type of menu (catalog, menu, button)
	Method     string    `json:"method" gorm:"size:20;index;"`                 // Http method of resource
	Path       string    `json:"path" gorm:"size:255;"`                        // Access path of menu
	Component  string    `json:"component" gorm:"size:255;"`                   // Component path of view
	Status     string    `json:"status" gorm:"size:20;index"`                  // Status of menu (enabled, disabled)
	Redirect   string    `json:"redirect" gorm:"size:255;not null;default:''"` // Redirect path of menu
	ParentID   string    `json:"parentId" gorm:"size:20;index;"`               // Parent ID (From Menu.ID)
	ParentPath string    `json:"parentPath" gorm:"size:255;index;"`            // Parent path (split by .)
	Children   *Menus    `json:"children" gorm:"-"`                            // Child menus
	Meta       MenuMeta  `json:"meta" gorm:"embedded"`                         // Meta of menu
	CreatedAt  time.Time `json:"createdAt" gorm:"index;"`                      // Create time
	UpdatedAt  time.Time `json:"updatedAt" gorm:"index;"`                      // Update time

	Roles Roles `json:"roles" gorm:"many2many:role_menus;"`
}

func (a Menu) TableName() string {
	return configs.C.FormatTableName("menu")
}

// Defining the slice of `Menu` struct.
type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	if a[i].Meta.Rank == a[j].Meta.Rank {
		return a[i].CreatedAt.Unix() > a[j].CreatedAt.Unix()
	}
	return a[i].Meta.Rank > a[j].Meta.Rank
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

type MenuMeta struct {
	Rank  int    `json:"order" gorm:"column:rank;index;"` // Rank for sorting (Order by desc)
	Title string `json:"title" gorm:"size:1024"`          // Menu title

	Properties map[string]any `json:"-" gorm:"type:text;serializer:json;default:'{}'"`
}

func (m *MenuMeta) UnmarshalJSON(data []byte) error {
	var properties map[string]any
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}

	if properties != nil {
		if v, ok := properties["order"]; ok {
			if vv, ok := v.(float64); ok {
				m.Rank = int(vv)
			}
			delete(properties, "order")
		}

		if v, ok := properties["title"]; ok {
			if vv, ok := v.(string); ok {
				m.Title = vv
			}
			delete(properties, "title")
		}
	}
	m.Properties = properties
	return nil
}

func (m MenuMeta) MarshalJSON() ([]byte, error) {

	props := m.Properties
	if props == nil {
		props = map[string]any{}
	}

	props["order"] = m.Rank
	props["title"] = m.Title

	return json.Marshal(props)
}
