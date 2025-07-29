package models

import (
	"time"

	"gin-admin/internal/configs"
)

// Role permissions for SYS
type MenuRole struct {
	ID        string    `json:"id" gorm:"size:20;primarykey"`                          // Unique ID
	RoleID    string    `json:"roleId" gorm:"size:20;uniqueIndex:idx_role_menu_index"` // From Role.ID
	MenuID    string    `json:"menuId" gorm:"size:20;uniqueIndex:idx_role_menu_index"` // From Menu.ID
	CreatedAt time.Time `json:"createdAt" gorm:"index;"`                               // Create time
	UpdatedAt time.Time `json:"updatedAt" gorm:"index;"`                               // Update time
}

func (a MenuRole) TableName() string {
	return configs.C.FormatTableName("role_menus")
}

// Defining the slice of `MenuRole` struct.
type MenuRoles []*MenuRole
