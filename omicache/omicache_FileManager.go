package cache

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FileManager struct {
	BaseDir string
}

// 初始化文件管理器
func NewFileManager(baseDir string) *FileManager {
	_ = os.MkdirAll(baseDir, 0755) // 确保目录存在
	return &FileManager{BaseDir: baseDir}
}

// 写文件
func (fm *FileManager) WriteFile(filename string, data []byte) error {
	path := filepath.Join(fm.BaseDir, fm.sanitizeFilename(filename))
	return os.WriteFile(path, data, 0644)
}

// 读文件
func (fm *FileManager) ReadFile(filename string) ([]byte, error) {
	path := filepath.Join(fm.BaseDir, fm.sanitizeFilename(filename))
	return os.ReadFile(path)
}

// 删除文件
func (fm *FileManager) DeleteFile(filename string) error {
	path := filepath.Join(fm.BaseDir, fm.sanitizeFilename(filename))
	return os.Remove(path)
}

// 转换 URL 路径为安全的文件名
func (fm *FileManager) sanitizeFilename(filename string) string {
	return strings.ReplaceAll(filename, "/", "@")
}

type FileInfo struct {
	Key  string // 文件名对应的缓存键
	Size int    // 文件大小（字节）
}

// 列出缓存目录中的所有文件及其大小
func (fm *FileManager) ListFiles() []FileInfo {
	var files []FileInfo

	// 遍历目录中的文件
	_ = filepath.Walk(fm.BaseDir, func(path string, info fs.FileInfo, err error) error {
		// 如果出错或是目录，则跳过
		if err != nil || info.IsDir() {
			return nil
		}

		// 转换文件名为缓存键
		relativePath, _ := filepath.Rel(fm.BaseDir, path)
		cacheKey := strings.ReplaceAll(relativePath, "@", "/")

		// 添加文件信息
		files = append(files, FileInfo{
			Key:  cacheKey,
			Size: int(info.Size()),
		})
		return nil
	})

	return files
}
