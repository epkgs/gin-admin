package oss

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LocalClientConfig struct {
	Domain     string
	RootPath   string
	BucketName string
	Prefix     string
}

var _ IClient = (*LocalClient)(nil)

// LocalObjectModel 用于存储本地对象的元数据
type LocalObjectModel struct {
	ID           uint              `gorm:"primaryKey;autoIncrement"`
	Key          string            `gorm:"uniqueIndex;size:512"`
	ETag         string            `gorm:"size:64"`
	Size         int64             `gorm:"index"`
	ContentType  string            `gorm:"size:128"`
	UserMetadata map[string]string `gorm:"type:text"`
	LastModified time.Time         `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m *LocalObjectModel) TableName() string {
	return "local_oss"
}

type LocalClient struct {
	config LocalClientConfig
	db     *gorm.DB
}

func NewLocalClient(config LocalClientConfig, db *gorm.DB) (*LocalClient, error) {
	// 确保根目录存在
	if err := os.MkdirAll(config.RootPath, 0755); err != nil {
		return nil, err
	}

	// 如果提供了数据库连接，则自动迁移表结构
	if db != nil {
		err := db.AutoMigrate(&LocalObjectModel{})
		if err != nil {
			return nil, err
		}
	}

	return &LocalClient{
		config: config,
		db:     db,
	}, nil
}

func (c *LocalClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.ReadSeeker, objectSize int64, options ...func(*PutObjectOptions)) (*PutObjectResult, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	opt := PutObjectOptions{
		ContentType:  "",
		UserMetadata: map[string]string{},
	}
	for _, f := range options {
		f(&opt)
	}

	objectName = formatObjectName(c.config.Prefix, objectName)

	// 构建完整路径
	fullPath := filepath.Join(c.config.RootPath, bucketName, objectName)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// 创建目标文件
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		etag    string
		written int64
	)

	// 只有当数据库存在时才计算ETag
	if c.db != nil {
		hash := md5.New()
		teeReader := io.TeeReader(reader, hash)
		// 直接复制数据到目标文件并同时计算MD5
		written, err = io.Copy(file, teeReader)
		if err != nil {
			return nil, err
		}

		// 计算ETag (MD5)
		etag = hex.EncodeToString(hash.Sum(nil))
	} else {
		// 直接复制数据
		written, err = io.Copy(file, reader)
		if err != nil {
			return nil, err
		}
	}

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// 确定MIME类型
	contentType := opt.ContentType
	if contentType == "" {
		// 根据文件扩展名推断MIME类型
		ext := filepath.Ext(fullPath)
		contentType = mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	// 如果有数据库连接，保存元数据
	if c.db != nil {
		model := &LocalObjectModel{
			Key:          objectName,
			ETag:         etag,
			Size:         info.Size(),
			ContentType:  contentType,
			UserMetadata: opt.UserMetadata,
			LastModified: info.ModTime(),
		}

		// 使用Upsert操作保存元数据
		err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var existing LocalObjectModel
			result := tx.Where("key = ?", objectName).First(&existing)
			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					// 创建新记录
					return tx.Create(model).Error
				}
				return result.Error
			}
			// 更新现有记录
			model.ID = existing.ID
			return tx.Save(model).Error
		})
		if err != nil {
			return nil, err
		}
	}

	return &PutObjectResult{
		URL:  c.config.Domain + "/" + filepath.ToSlash(objectName),
		Key:  objectName,
		ETag: etag,
		Size: written,
	}, nil
}

func (c *LocalClient) GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)

	// 构建完整路径
	fullPath := filepath.Join(c.config.RootPath, bucketName, objectName)

	// 打开文件
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, err
	}

	return file, nil
}

func (c *LocalClient) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)

	// 构建完整路径
	fullPath := filepath.Join(c.config.RootPath, bucketName, objectName)

	// 删除文件
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 如果有数据库连接，删除元数据
	if c.db != nil {
		err = c.db.WithContext(ctx).Where("key = ?", objectName).Delete(&LocalObjectModel{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *LocalClient) RemoveObjectByURL(ctx context.Context, urlStr string) error {
	prefix := c.config.Domain
	if !strings.HasPrefix(urlStr, prefix) {
		return nil
	}

	objectName := strings.TrimPrefix(urlStr, prefix)
	objectName = strings.TrimPrefix(objectName, "/") // 移除开头的斜杠

	return c.RemoveObject(ctx, "", objectName)
}

func (c *LocalClient) StatObject(ctx context.Context, bucketName, objectName string) (*ObjectStat, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)

	// 如果有数据库连接，优先从数据库获取元数据
	if c.db != nil {
		var model LocalObjectModel
		result := c.db.WithContext(ctx).Where("key = ?", objectName).First(&model)
		if result.Error == nil {
			return &ObjectStat{
				Key:          model.Key,
				ETag:         model.ETag,
				Size:         model.Size,
				ContentType:  model.ContentType,
				UserMetadata: model.UserMetadata,
				LastModified: model.LastModified,
			}, nil
		}
		// 如果数据库查询失败，继续使用文件系统方式
	}

	// 构建完整路径
	fullPath := filepath.Join(c.config.RootPath, bucketName, objectName)

	// 获取文件信息
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, err
	}

	// 确定MIME类型
	contentType := "application/octet-stream"
	ext := filepath.Ext(fullPath)
	if mimeTypes, err := mime.ExtensionsByType(contentType); err == nil && len(mimeTypes) > 0 {
		contentType = mime.TypeByExtension(ext)
	}

	return &ObjectStat{
		Key:          objectName,
		Size:         info.Size(),
		ETag:         "", // 本地文件系统本身不支持ETag，但如果有数据库会从数据库获取
		LastModified: info.ModTime(),
		ContentType:  contentType,
		UserMetadata: make(map[string]string),
	}, nil
}

func (c *LocalClient) StatObjectByURL(ctx context.Context, urlStr string) (*ObjectStat, error) {
	prefix := c.config.Domain
	if !strings.HasPrefix(urlStr, prefix) {
		return nil, nil
	}

	objectName := strings.TrimPrefix(urlStr, prefix)
	// 移除开头的斜杠
	objectName = strings.TrimPrefix(objectName, "/")

	return c.StatObject(ctx, "", objectName)
}
