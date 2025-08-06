package repositories

import "gin-admin/pkg/gormx"

type base[T gormx.Entity] = gormx.Repository[T]
