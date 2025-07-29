package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gin-admin/internal/configs"
)

const (
	UserStatus_Activated = "activated"
	UserStatus_Freezed   = "freezed"
)

// User management for SYS
type User struct {
	ID          string    `json:"id" gorm:"size:20;primarykey;"`                                                       // Unique ID
	Username    string    `json:"username" gorm:"size:64;index"`                                                       // Username for login
	Password    string    `json:"-" gorm:"size:64;"`                                                                   // Password for login (encrypted)
	NickName    string    `json:"nickName" gorm:"size:64;index"`                                                       // Name of user
	RealName    string    `json:"realName" gorm:"size:64;"`                                                            // Real name of user
	Wechat      string    `json:"wechat" gorm:"size:64;"`                                                              // Wechat account
	Phone       string    `json:"phone" gorm:"size:32;"`                                                               // Phone number of user
	Email       string    `json:"email" gorm:"size:128;"`                                                              // Email of user
	Status      string    `json:"status" gorm:"size:20;index"`                                                         // Status of user (activated, freezed)
	Description string    `json:"description" gorm:"size:1024"`                                                        // Details about user
	Avatar      string    `json:"avatar" gorm:"not null;default:'';comment:Avatar URL"`                                // Avatar URL
	Fingers     Fingers   `json:"-" gorm:"type:string;serializer:json;not null;default:'[]';comment:Fingerprint list"` // Frontend fingerprints
	CreatedAt   time.Time `json:"createdAt" gorm:"index;"`                                                             // Create time
	UpdatedAt   time.Time `json:"updatedAt" gorm:"index;"`                                                             // Update time

	Roles Roles `json:"roles" gorm:"many2many:user_roles;"` // Roles of user
}

func (a User) TableName() string {
	return configs.C.FormatTableName("user")
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

type Fingers [2]string

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据并解析为 []string
func (f *Fingers) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("assert finger attribute to string failed")
	}

	return json.Unmarshal([]byte(str), f)
}

// Value 实现 driver.Valuer 接口，用于将 []string 序列化为 JSON 字符串存储到数据库
func (f Fingers) Value() (driver.Value, error) {
	byts, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}
