package model

import (
	"encoding/json"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const (
	RoleStatus_Enabled  = "enabled"  // Enabled
	RoleStatus_Disabled = "disabled" // Disabled

	RoleResultType_Select = "select" // Select
)

func init() {
	AddMigrationModel(Role{})
}

// Role management
type Role struct {
	ID          string                `json:"id" gorm:"size:20;primarykey;"`              // Unique ID
	Code        string                `json:"code" gorm:"size:32;index;"`                 // Code of role (unique)
	Name        string                `json:"name" gorm:"size:128;index"`                 // Display name of role
	Description string                `json:"description" gorm:"size:1024"`               // Details about role
	Rank        int                   `json:"rank" gorm:"index"`                          // Rank for sorting
	Status      string                `json:"status" gorm:"size:20;index"`                // Status of role (disabled, enabled)
	CreatedAt   time.Time             `json:"createdAt" gorm:"index;"`                    // Create time
	UpdatedAt   time.Time             `json:"updatedAt" gorm:"index;"`                    // Update time
	DeletedAt   soft_delete.DeletedAt `json:"-" gorm:"index;default:0;softDelete:milli;"` // Delete time

	Menus Menus `json:"menus" gorm:"many2many:r_role_menus;"`
	Users Users `json:"users" gorm:"many2many:r_role_users;"`
}

func (a Role) TableName() string {
	return "role"
}

func (a *Role) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = xid.New().String()
	}
	return nil
}

// Defining the slice of `Role` struct.
type Roles []*Role

func (m Roles) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("[]"), nil
	}

	type M Roles
	copy := M(m)

	return json.Marshal(copy)
}
