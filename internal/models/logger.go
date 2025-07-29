package models

import (
	"gin-admin/internal/configs"
	"gin-admin/pkg/logging"
)

// Logger management
type Logger struct {
	logging.Logger

	NickName string `json:"nickName" gorm:"<-:false;-:migration;"` // From User.NickName
	Username string `json:"username" gorm:"<-:false;-:migration;"` // From User.Name
}

func (a Logger) TableName() string {
	return configs.C.FormatTableName("logger")
}

// Defining the slice of `Logger` struct.
type Loggers []*Logger
