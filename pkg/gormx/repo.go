package gormx

import (
	"context"

	"gorm.io/gorm"
)

// Entity 实体接口
// 所有可以被仓库管理的实体都应该实现这个接口
type Entity interface {
	TableName() string
}

// Repository 数据库操作接口
// 提供基本的CRUD操作
type Repository[T Entity] interface {
	// Create 创建实体
	Create(ctx context.Context, entity *T, opts ...Option) error

	// CreateBatch 批量创建实体
	CreateBatch(ctx context.Context, entities []*T, batchSize int, opts ...Option) error

	// Get 根据ID或条件获取单个实体
	// id: 实体ID，如果为nil，则使用condition和args
	// opts: 查询选项
	Get(ctx context.Context, id any, opts ...Option) (*T, error)

	// First 根据条件获取第一个实体
	// opts: 查询选项
	First(ctx context.Context, opts ...Option) (*T, error)

	// Update 更新实体
	Update(ctx context.Context, entity *T, opts ...Option) error

	// Delete 删除实体
	Delete(ctx context.Context, id any, opts ...Option) error

	// Delete 删除实体
	DeleteBatch(ctx context.Context, opts ...Option) error

	// Find 查询实体列表
	// opts: 查询选项
	Find(ctx context.Context, opts ...Option) ([]*T, error)

	// Count 获取实体总数
	// opts: 查询选项
	Count(ctx context.Context, opts ...Option) (int64, error)

	// Exists 检查实体是否存在
	Exists(ctx context.Context, opts ...Option) (bool, error)

	// Transaction 在事务中执行函数
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	// WithTx 使用事务
	WithTx(tx *gorm.DB) Repository[T]

	// DB 实例
	DB() *gorm.DB
}

// GenericRepo 通用仓库实现
type GenericRepo[T Entity] struct {
	db *gorm.DB
}

// NewGenericRepo 创建通用仓库
func NewGenericRepo[T Entity](db *gorm.DB) *GenericRepo[T] {
	return &GenericRepo[T]{
		db: db,
	}
}

// Create 创建实体
func (r *GenericRepo[T]) Create(ctx context.Context, entity *T, opts ...Option) error {
	query := Apply(r.db.WithContext(ctx), opts...)
	return query.Create(entity).Error
}

func (r *GenericRepo[T]) CreateBatch(ctx context.Context, entities []*T, batchSize int, opts ...Option) error {
	query := Apply(r.db.WithContext(ctx), opts...)
	return query.CreateInBatches(entities, batchSize).Error
}

// Get 根据ID或条件获取单个实体
func (r *GenericRepo[T]) Get(ctx context.Context, id any, opts ...Option) (*T, error) {
	var entity T

	// 创建查询并应用选项
	query := Apply(r.db.WithContext(ctx), opts...)

	if id != nil {
		query = query.Where("id = ?", id)
	}

	if len(query.Statement.Clauses) == 0 {
		// 没有ID和条件，返回错误
		return nil, gorm.ErrMissingWhereClause
	}

	// 根据条件查询
	err := query.First(&entity).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &entity, err
}

// First 查找第一个
func (r *GenericRepo[T]) First(ctx context.Context, opts ...Option) (*T, error) {
	var entity T

	// 创建查询并应用选项
	query := Apply(r.db.WithContext(ctx), opts...)

	if len(query.Statement.Clauses) == 0 {
		// 没有查询条件，返回错误
		return nil, gorm.ErrMissingWhereClause
	}

	// 根据条件查询
	err := query.First(&entity).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &entity, err
}

// Update 更新实体
func (r *GenericRepo[T]) Update(ctx context.Context, entity *T, opts ...Option) error {
	query := Apply(r.db.WithContext(ctx), opts...)
	return query.Updates(entity).Error
}

// Delete 删除实体
func (r *GenericRepo[T]) Delete(ctx context.Context, id any, opts ...Option) error {
	var entity T
	query := Apply(r.db.WithContext(ctx), opts...)
	return query.Delete(&entity, id).Error
}

// DeleteBatch 删除实体
func (r *GenericRepo[T]) DeleteBatch(ctx context.Context, opts ...Option) error {
	var entity T
	query := Apply(r.db.WithContext(ctx), opts...)

	if len(query.Statement.Clauses) == 0 {
		// 没有查询条件，返回错误
		return gorm.ErrMissingWhereClause
	}

	return query.Delete(&entity).Error
}

// Find 查询实体列表
func (r *GenericRepo[T]) Find(ctx context.Context, opts ...Option) ([]*T, error) {
	var entities []*T

	// 创建查询
	query := r.db.WithContext(ctx)

	// 应用查询选项
	query = Apply(query, opts...)

	// 执行查询
	err := query.Find(&entities).Error

	if err == gorm.ErrRecordNotFound {
		return make([]*T, 0), nil
	}

	return entities, err
}

// Count 获取实体总数
func (r *GenericRepo[T]) Count(ctx context.Context, opts ...Option) (int64, error) {
	var count int64
	var entity T

	// 创建查询
	query := r.db.WithContext(ctx).Model(&entity)

	// 应用查询选项
	query = Apply(query, opts...)

	// 执行查询
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Exists 判断实体是否存在
func (r *GenericRepo[T]) Exists(ctx context.Context, opts ...Option) (bool, error) {
	count, err := r.Count(ctx, opts...)
	return count > 0, err
}

// Transaction 在事务中执行函数
func (r *GenericRepo[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// WithTx 使用事务
func (r *GenericRepo[T]) WithTx(tx *gorm.DB) Repository[T] {
	return &GenericRepo[T]{
		db: tx,
	}
}

func (r *GenericRepo[T]) DB() *gorm.DB {
	return r.db
}
