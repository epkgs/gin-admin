package oss

import (
	"context"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClientConfig struct {
	Domain          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Prefix          string
}

var _ IClient = (*MinioClient)(nil)

type MinioClient struct {
	config MinioClientConfig
	client *minio.Client
}

func NewMinioClient(config MinioClientConfig) (*MinioClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if exists, err := client.BucketExists(ctx, config.BucketName); err != nil {
		return nil, err
	} else if !exists {
		if err := client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		config: config,
		client: client,
	}, nil
}

func (c *MinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.ReadSeeker, objectSize int64, options ...func(*PutObjectOptions)) (*PutObjectResult, error) {
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
	output, err := c.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType:  opt.ContentType,
		UserMetadata: opt.UserMetadata,
	})
	if err != nil {
		return nil, err
	}

	return &PutObjectResult{
		URL:  c.config.Domain + "/" + objectName,
		Key:  output.Key,
		ETag: output.ETag,
		Size: output.Size,
	}, nil
}

func (c *MinioClient) GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	return c.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (c *MinioClient) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	return c.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (c *MinioClient) GetObjectByURL(ctx context.Context, urlStr string) (io.ReadCloser, error) {
	bucketName, objectName := c.splitBucketAndObject(urlStr)
	return c.GetObject(ctx, bucketName, objectName)
}

func (c *MinioClient) RemoveObjectByURL(ctx context.Context, urlStr string) error {
	bucketName, objectName := c.splitBucketAndObject(urlStr)
	return c.RemoveObject(ctx, bucketName, objectName)
}

func (c *MinioClient) StatObjectByURL(ctx context.Context, urlStr string) (*ObjectStat, error) {
	bucketName, objectName := c.splitBucketAndObject(urlStr)
	return c.StatObject(ctx, bucketName, objectName)
}

func (c *MinioClient) StatObject(ctx context.Context, bucketName, objectName string) (*ObjectStat, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	info, err := c.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &ObjectStat{
		Key:          info.Key,
		Size:         info.Size,
		ETag:         info.ETag,
		ContentType:  info.ContentType,
		UserMetadata: info.UserMetadata,
	}, nil
}

func (c *MinioClient) splitBucketAndObject(urlStr string) (bucketName, objectName string) {
	prefix := c.config.Domain + "/"
	if !strings.HasPrefix(urlStr, prefix) {
		return "", ""
	}

	objectName = strings.TrimPrefix(urlStr, prefix)
	return c.config.BucketName, objectName
}
