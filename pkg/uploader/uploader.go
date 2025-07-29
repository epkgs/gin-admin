package uploader

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type Uploader struct {
	option *Option
}

type Option struct {
	UploadPath string
	UseDateDir bool
}

type FileInfo struct {
	Path string `json:"path"` // 相对路径，包含文件名和后缀 ( 统一为 unix 风格 '/' )
	Name string `json:"name"` // 仅文件名，不包含后缀
	Ext  string `json:"ext"`  // 文件后缀，包含 "."
	Mime string `json:"mime"` // MIME类型
	Size int64  `json:"size"` // 文件大小
}

func New(opts ...func(opt *Option)) *Uploader {
	opt := &Option{
		UploadPath: "uploads",
		UseDateDir: true,
	}

	for _, fn := range opts {
		fn(opt)
	}

	return &Uploader{opt}
}

func (up *Uploader) Upload(ctx context.Context, header *multipart.FileHeader) (*FileInfo, error) {
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := &FileInfo{}
	info.Size = header.Size
	info.Mime = header.Header.Get("Content-Type")
	info.Ext = getExt(header.Filename, info.Mime)
	info.Name = header.Filename[:len(header.Filename)-len(info.Ext)]

	// 创建文件
	out, err := up.createFile(info)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// 将上传的文件内容复制到本地文件
	_, err = io.Copy(out, file)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (up *Uploader) Delete(ctx context.Context, path string) error {
	return os.Remove(path)
}

func (up *Uploader) Open(ctx context.Context, path string) (file *os.File, err error) {
	return os.Open(path)
}

func (up *Uploader) Exists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (up *Uploader) createFile(info *FileInfo) (*os.File, error) {

	uploadPath := filepath.ToSlash(up.option.UploadPath)
	var buildPath = func() {
		if up.option.UseDateDir {
			info.Path = filepath.Join(uploadPath, time.Now().Format("20060102"), info.Name+info.Ext)
		} else {
			info.Path = filepath.Join(uploadPath, info.Name+info.Ext)
		}
	}

	buildPath()

	idx := 0
	for fileExists(info.Path) {
		idx++

		if idx > 999 {
			return nil, fmt.Errorf("文件名已存在，且已超过最大重命名次数")
		}

		info.Name = fmt.Sprintf("%s_%03d", info.Name, idx)

		buildPath()
	}

	dir, _ := filepath.Split(info.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err // 如果创建目录失败，则返回错误
	}

	return os.Create(info.Path) // 创建文件
}

func getExt(fileName, contentType string) string {
	if ext := filepath.Ext(fileName); ext != "" {
		return ext
	}

	// 如果需要更详细的MIME类型信息，可以使用mime包进一步解析
	// 例如，处理 application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
	extension, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return ""
	}

	return extension[0] // 选第一个包含点的后缀名，例如 ".xlsx"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// 其他错误
	return false
}
