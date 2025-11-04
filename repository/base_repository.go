package repository

import (
	"context"

	"gorm.io/gorm"
)

type dbTxKey struct{}

type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建 BaseRepository
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return BaseRepository[T]{db: db}
}

// Transaction 开启事务
func (a *BaseRepository[T]) Transaction(c context.Context, handler func(tx context.Context) error) error {
	db := a.db
	return db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		return handler(context.WithValue(c, dbTxKey{}, tx))
	})
}

// GetDb 获取事务的db连接
func (a *BaseRepository[T]) GetDb(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(dbTxKey{}).(*gorm.DB)
	if !ok {
		db = a.db
		return db.WithContext(ctx)
	}
	return db
}

// GetById 根据ID查询
func (a *BaseRepository[T]) GetById(ctx context.Context, id int64) (T, error) {
	var domain T
	err := a.GetDb(ctx).Limit(1).Find(&domain, id).Error
	return domain, err
}

// Insert 新增
func (a *BaseRepository[T]) Insert(ctx context.Context, entity T) (T, error) {
	ret := a.GetDb(ctx).Create(&entity)
	return entity, ret.Error
}

// Update 更新
func (a *BaseRepository[T]) Update(ctx context.Context, entity T) (T, error) {
	db := a.GetDb(ctx).Model(&entity).Updates(entity)
	return entity, db.Error
}

// DeleteById 根据ID删除
func (a *BaseRepository[T]) DeleteById(ctx context.Context, id int64) error {
	db := a.GetDb(ctx).Delete(new(T), id)
	return db.Error
}

// ListByIds 根据ID集合查询
func (a *BaseRepository[T]) ListByIds(ctx context.Context, ids []int64) ([]T, error) {
	var list []T
	db := a.GetDb(ctx).Find(&list, ids)
	return list, db.Error
}
