package models

import (
	"time"

	"gin-admin/internal/configs"
)

// User roles association
type UserRole struct {
	ID        string    `json:"id" gorm:"size:20;primarykey"`          // Unique ID
	UserID    string    `json:"userId" gorm:"size:20;index"`           // From User.ID
	RoleID    string    `json:"roleId" gorm:"size:20;index"`           // From Role.ID
	CreatedAt time.Time `json:"createdAt" gorm:"index;"`               // Create time
	UpdatedAt time.Time `json:"updatedAt" gorm:"index;"`               // Update time
	RoleName  string    `json:"roleName" gorm:"<-:false;-:migration;"` // From Role.Name
}

func (a UserRole) TableName() string {
	return configs.C.FormatTableName("user_roles")
}

// Defining the slice of `UserRole` struct.
type UserRoles []*UserRole

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, userRole := range a {
		m[userRole.UserID] = append(m[userRole.UserID], userRole)
	}
	return m
}

func (a UserRoles) ToRoleIDs() []string {
	var ids []string
	for _, item := range a {
		ids = append(ids, item.RoleID)
	}
	return ids
}
