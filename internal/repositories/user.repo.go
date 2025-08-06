package repositories

import (
	"context"

	"gin-admin/internal/models"
	"gin-admin/pkg/gormx"

	"gorm.io/gorm"
)

// User management for SYS
type User struct {
	base[models.User]
}

func NewUser(db *gorm.DB) *User {
	return &User{
		base: gormx.NewGenericRepo[models.User](db),
	}
}

func (a *User) GetByUsername(ctx context.Context, username string, opts ...gormx.Option) (*models.User, error) {

	return a.First(ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("username = ?", username)
		return gormx.Apply(db, opts...)
	})
}

func (a *User) ExistsUsername(ctx context.Context, username string) (bool, error) {
	return a.base.Exists(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("username=?", username)
	})
}

func (a *User) UpdatePassword(ctx context.Context, id string, password string) error {
	user := &models.User{
		ID:       id,
		Password: password,
	}
	return a.Update(ctx, user, gormx.WithSelect("Password"))
}
