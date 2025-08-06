package repositories

import (
	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Logger management
type Logger struct {
	base[models.Logger]
}

func NewLogger(db *gorm.DB) *Logger {
	return &Logger{
		base: gormx.NewGenericRepo[models.Logger](db),
	}
}
