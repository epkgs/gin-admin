package po

import (
	"time"
)

// User roles association
type RoleUser struct {
	ID        string    `json:"id" gorm:"size:20;primarykey"`          // Unique ID
	UserID    string    `json:"userId" gorm:"size:20;index"`           // From User.ID
	RoleID    string    `json:"roleId" gorm:"size:20;index"`           // From Role.ID
	CreatedAt time.Time `json:"createdAt" gorm:"index;"`               // Create time
	UpdatedAt time.Time `json:"updatedAt" gorm:"index;"`               // Update time
	RoleName  string    `json:"roleName" gorm:"<-:false;-:migration;"` // From Role.Name
}

func (a RoleUser) TableName() string {
	return "r_role_users"
}

// Defining the slice of `RoleUser` struct.
type RoleUsers []*RoleUser

func (a RoleUsers) ToUserIDMap() map[string]RoleUsers {
	m := make(map[string]RoleUsers)
	for _, userRole := range a {
		m[userRole.UserID] = append(m[userRole.UserID], userRole)
	}
	return m
}

func (a RoleUsers) ToRoleIDs() []string {
	var ids []string
	for _, item := range a {
		ids = append(ids, item.RoleID)
	}
	return ids
}
