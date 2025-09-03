package cachex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
)

// BadgerDB 缓存实现
type badgerCache struct {
	db *badger.DB
}

type BadgerConfig struct {
	Path string // 数据库存储路径
}

func NewBadgerCache(config *BadgerConfig) (Cacher, error) {
	// 设置 BadgerDB 选项
	opts := badger.DefaultOptions(config.Path)
	opts.Logger = nil // 禁用日志，可根据需要启用

	// 打开数据库
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger db: %w", err)
	}

	return &badgerCache{db: db}, nil
}

func (b *badgerCache) Get(ctx context.Context, key string) (value []byte, err error) {
	err = b.db.View(func(txn *badger.Txn) error {
		var item *badger.Item
		item, err = txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}

		return item.Value(func(val []byte) error {
			value = val
			return nil
		})
	})

	return
}

func (b *badgerCache) GetObject(ctx context.Context, key string, value interface{}) error {
	raw, err := b.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, value)
}

func (b *badgerCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value)
		if expiration > 0 {
			e = e.WithTTL(expiration)
		}
		return txn.SetEntry(e)
	})
}

func (b *badgerCache) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return b.Set(ctx, key, data, expiration)
}

func (b *badgerCache) Delete(ctx context.Context, key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *badgerCache) Exists(ctx context.Context, key string) (bool, error) {
	exists := false
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		exists = true
		return nil
	})
	return exists, err
}

func (b *badgerCache) ForEach(namespace string, fn func(key string, raw []byte) bool) error {

	errCallbackInterrupted := errors.New("badger: iteration interrupted")

	return b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true

		it := txn.NewIterator(opts)
		defer it.Close()

		// 构建前缀
		var prefix []byte
		if namespace != "" {
			prefix = []byte(namespace + Delimiter)
		}

		// 遍历所有匹配的键值对
		for it.Seek(prefix); ; it.Next() {
			// 检查是否还有有效数据
			if len(prefix) > 0 && !it.ValidForPrefix(prefix) {
				break
			} else if !it.Valid() {
				break
			}

			item := it.Item()
			key := item.Key()

			// 提取子键（去掉命名空间前缀）
			var subKey string
			if namespace == "" {
				subKey = string(key)
			} else {
				subKey = string(key[len(prefix):])
			}

			// 处理值并调用回调函数
			err := item.Value(func(val []byte) error {
				if !fn(subKey, val) {
					return errCallbackInterrupted // 明确表示迭代被中断
				}
				return nil
			})

			// 如果迭代被中断，直接返回nil而不是错误
			if errors.Is(err, errCallbackInterrupted) {
				return nil
			}

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (b *badgerCache) Close() error {
	return b.db.Close()
}
