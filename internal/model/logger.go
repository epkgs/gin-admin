package model

import (
	"gin-admin/pkg/logger"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

func init() {
	AddMigrationModel(Logger{})
}

// Logger management
type Logger struct {
	logger.Logger

	CreatedAt time.Time `json:"createdAt" gorm:"index;"` // Create time
	UpdatedAt time.Time `json:"updatedAt" gorm:"index;"` // Update time

	NickName string `json:"nickName" gorm:"<-:false;-:migration;"` // From User.NickName
	Username string `json:"username" gorm:"<-:false;-:migration;"` // From User.Name
}

func (a Logger) TableName() string {
	return "logger"
}

func (a *Logger) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = xid.New().String()
	}
	return nil
}

// Defining the slice of `Logger` struct.
type Loggers []*Logger
