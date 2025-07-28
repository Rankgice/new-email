package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageConfig 存储配置
type StorageConfig struct {
	Type      string   `json:"type"`      // local, oss, s3
	BasePath  string   `json:"basePath"`  // 本地存储基础路径
	MaxSize   int64    `json:"maxSize"`   // 最大文件大小（字节）
	AllowExts []string `json:"allowExts"` // 允许的文件扩展名
	CDNDomain string   `json:"cdnDomain"` // CDN域名
}

// StorageService 存储服务
type StorageService struct {
	config StorageConfig
}

// NewStorageService 创建存储服务
func NewStorageService(config StorageConfig) *StorageService {
	return &StorageService{
		config: config,
	}
}

// FileInfo 文件信息
type FileInfo struct {
	Filename    string    `json:"filename"`    // 原始文件名
	StorageName string    `json:"storageName"` // 存储文件名
	Path        string    `json:"path"`        // 存储路径
	URL         string    `json:"url"`         // 访问URL
	Size        int64     `json:"size"`        // 文件大小
	ContentType string    `json:"contentType"` // 内容类型
	MD5         string    `json:"md5"`         // MD5值
	UploadTime  time.Time `json:"uploadTime"`  // 上传时间
}

// UploadFile 上传文件
func (s *StorageService) UploadFile(file multipart.File, header *multipart.FileHeader, category string) (*FileInfo, error) {
	// 检查文件大小
	if header.Size > s.config.MaxSize {
		return nil, fmt.Errorf("文件大小超过限制: %d > %d", header.Size, s.config.MaxSize)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !s.isAllowedExt(ext) {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 生成存储文件名
	storageName := s.generateStorageName(header.Filename)

	// 构建存储路径
	storagePath := s.buildStoragePath(category, storageName)

	// 确保目录存在
	if err := s.ensureDir(filepath.Dir(storagePath)); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存文件
	if err := s.saveFile(file, storagePath); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 计算MD5
	md5Hash, err := s.calculateMD5(storagePath)
	if err != nil {
		return nil, fmt.Errorf("计算MD5失败: %v", err)
	}

	// 构建访问URL
	url := s.buildURL(category, storageName)

	return &FileInfo{
		Filename:    header.Filename,
		StorageName: storageName,
		Path:        storagePath,
		URL:         url,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
		MD5:         md5Hash,
		UploadTime:  time.Now(),
	}, nil
}

// UploadBytes 上传字节数据
func (s *StorageService) UploadBytes(data []byte, filename, category string) (*FileInfo, error) {
	// 检查文件大小
	if int64(len(data)) > s.config.MaxSize {
		return nil, fmt.Errorf("文件大小超过限制: %d > %d", len(data), s.config.MaxSize)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if !s.isAllowedExt(ext) {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 生成存储文件名
	storageName := s.generateStorageName(filename)

	// 构建存储路径
	storagePath := s.buildStoragePath(category, storageName)

	// 确保目录存在
	if err := s.ensureDir(filepath.Dir(storagePath)); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存文件
	if err := os.WriteFile(storagePath, data, 0644); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 计算MD5
	md5Hash := fmt.Sprintf("%x", md5.Sum(data))

	// 构建访问URL
	url := s.buildURL(category, storageName)

	return &FileInfo{
		Filename:    filename,
		StorageName: storageName,
		Path:        storagePath,
		URL:         url,
		Size:        int64(len(data)),
		MD5:         md5Hash,
		UploadTime:  time.Now(),
	}, nil
}

// DeleteFile 删除文件
func (s *StorageService) DeleteFile(path string) error {
	fullPath := filepath.Join(s.config.BasePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}

// FileExists 检查文件是否存在
func (s *StorageService) FileExists(path string) bool {
	fullPath := filepath.Join(s.config.BasePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetFileInfo 获取文件信息
func (s *StorageService) GetFileInfo(path string) (*FileInfo, error) {
	fullPath := filepath.Join(s.config.BasePath, path)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 计算MD5
	md5Hash, err := s.calculateMD5(fullPath)
	if err != nil {
		return nil, fmt.Errorf("计算MD5失败: %v", err)
	}

	return &FileInfo{
		Filename:    stat.Name(),
		StorageName: stat.Name(),
		Path:        fullPath,
		Size:        stat.Size(),
		MD5:         md5Hash,
		UploadTime:  stat.ModTime(),
	}, nil
}

// isAllowedExt 检查是否为允许的文件扩展名
func (s *StorageService) isAllowedExt(ext string) bool {
	if len(s.config.AllowExts) == 0 {
		return true // 如果没有限制，则允许所有类型
	}

	for _, allowExt := range s.config.AllowExts {
		if strings.ToLower(allowExt) == ext {
			return true
		}
	}
	return false
}

// generateStorageName 生成存储文件名
func (s *StorageService) generateStorageName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d%s", timestamp, ext)
}

// buildStoragePath 构建存储路径
func (s *StorageService) buildStoragePath(category, filename string) string {
	// 按日期分目录
	date := time.Now().Format("2006/01/02")
	return filepath.Join(s.config.BasePath, category, date, filename)
}

// buildURL 构建访问URL
func (s *StorageService) buildURL(category, filename string) string {
	date := time.Now().Format("2006/01/02")
	relativePath := fmt.Sprintf("%s/%s/%s", category, date, filename)

	if s.config.CDNDomain != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(s.config.CDNDomain, "/"), relativePath)
	}

	return fmt.Sprintf("/uploads/%s", relativePath)
}

// ensureDir 确保目录存在
func (s *StorageService) ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// saveFile 保存文件
func (s *StorageService) saveFile(src multipart.File, dst string) error {
	// 重置文件指针
	src.Seek(0, 0)

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, src)
	return err
}

// calculateMD5 计算文件MD5
func (s *StorageService) calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetStorageStats 获取存储统计信息
func (s *StorageService) GetStorageStats() (map[string]interface{}, error) {
	var totalSize int64
	var fileCount int64

	err := filepath.Walk(s.config.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略错误，继续遍历
		}

		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalSize": totalSize,
		"fileCount": fileCount,
		"basePath":  s.config.BasePath,
		"maxSize":   s.config.MaxSize,
		"allowExts": s.config.AllowExts,
	}, nil
}

// CleanupOldFiles 清理旧文件
func (s *StorageService) CleanupOldFiles(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)

	return filepath.Walk(s.config.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略错误，继续遍历
		}

		if !info.IsDir() && info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("删除文件失败 %s: %v", path, err)
			}
		}

		return nil
	})
}
