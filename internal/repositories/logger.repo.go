package repositories

import (
	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// Logger management
type Logger struct {
	gormx.Repository[models.Logger]
}

func NewLogger(db *gorm.DB) *Logger {
	return &Logger{
		Repository: gormx.NewGenericRepo[models.Logger](db),
	}
}
