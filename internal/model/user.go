package model

import (
	"encoding/json"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const (
	UserStatus_Activated = "activated"
	UserStatus_Freezed   = "freezed"
)

func init() {
	AddMigrationModel(User{})
}

// User management for SYS
type User struct {
	ID          string                `json:"id" gorm:"size:20;primarykey;"`                        // Unique ID
	Username    string                `json:"username" gorm:"size:64;index"`                        // Username for login
	Password    string                `json:"-" gorm:"size:64;"`                                    // Password for login (encrypted)
	NickName    string                `json:"nickName" gorm:"size:64;index"`                        // Name of user
	RealName    string                `json:"realName" gorm:"size:64;"`                             // Real name of user
	Wechat      string                `json:"wechat" gorm:"size:64;"`                               // Wechat account
	Phone       string                `json:"phone" gorm:"size:32;"`                                // Phone number of user
	Email       string                `json:"email" gorm:"size:128;"`                               // Email of user
	Status      string                `json:"status" gorm:"size:20;index"`                          // Status of user (activated, freezed)
	Description string                `json:"description" gorm:"size:1024"`                         // Details about user
	Avatar      string                `json:"avatar" gorm:"not null;default:'';comment:Avatar URL"` // Avatar URL
	CreatedAt   time.Time             `json:"createdAt" gorm:"index;"`                              // Create time
	UpdatedAt   time.Time             `json:"updatedAt" gorm:"index;"`                              // Update time
	DeletedAt   soft_delete.DeletedAt `json:"-" gorm:"index;default:0;softDelete:milli;"`           // Delete time

	Roles Roles `json:"roles" gorm:"many2many:r_role_users;"` // Roles of user
}

func (a User) TableName() string {
	return "user"
}

func (a *User) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = xid.New().String()
	}
	return nil
}

// Defining the slice of `User` struct.
type Users []*User

func (a Users) ToIDs() []string {
	var ids []string
	for _, item := range a {
		ids = append(ids, item.ID)
	}
	return ids
}

func (m Users) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("[]"), nil
	}

	type M Users
	copy := M(m)

	return json.Marshal(copy)
}
