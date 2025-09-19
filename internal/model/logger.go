package model

import (
	"gin-admin/pkg/logger"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

func init() {
	AddMigrationModel(Logger{})
}

// Logger management
type Logger struct {
	logger.LoggerModel

	NickName string `json:"nickName" gorm:"<-:false;-:migration;"` // From User.NickName
	Username string `json:"username" gorm:"<-:false;-:migration;"` // From User.Name
}

func (a Logger) TableName() string {
	return logger.TableName
}

func (a *Logger) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = xid.New().String()
	}
	return nil
}

// Defining the slice of `Logger` struct.
type Loggers []*Logger
