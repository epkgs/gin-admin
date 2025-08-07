package oss

import (
	"context"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type LocalClientConfig struct {
	Domain     string
	RootPath   string
	BucketName string
	Prefix     string
}

var _ IClient = (*LocalClient)(nil)

type LocalClient struct {
	config LocalClientConfig
}

func NewLocalClient(config LocalClientConfig) (*LocalClient, error) {
	// 确保根目录存在
	if err := os.MkdirAll(config.RootPath, 0755); err != nil {
		return nil, err
	}

	return &LocalClient{
		config: config,
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

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 复制数据
	written, err := io.Copy(file, reader)
	if err != nil {
		return nil, err
	}

	// 获取文件信息
	_, err = file.Stat()
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

	return &PutObjectResult{
		URL:  c.config.Domain + "/" + filepath.ToSlash(objectName),
		Key:  objectName,
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
	return os.Remove(fullPath)
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
		ETag:         "", // 本地文件系统不支持ETag
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
