package web

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Web struct {
	SourcePath     string
	IndexPath      string
	EmbeddedSource *embed.FS
}

func NewWeb(sourcePath, indexPath string, embeddedSource *embed.FS) *Web {
	return &Web{
		SourcePath:     sourcePath,
		IndexPath:      indexPath,
		EmbeddedSource: embeddedSource,
	}
}

func (web *Web) GenerateTemplate() {
	copyEmbeddedFiles(web.SourcePath)
}

func (web *Web) ServeFile(w http.ResponseWriter, r *http.Request) error {
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = web.IndexPath
	}
	filePath = web.SourcePath + filePath

	var data []byte
	var err error
	if web.EmbeddedSource != nil {
		data, err = web.EmbeddedSource.ReadFile(filePath)
	} else {
		data, err = os.ReadFile(filePath)
	}

	// 检查是否读取失败
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return err
	}

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))

	// 设置Content-Type
	contentType := mimeByExtension(ext)
	w.Header().Set("Content-Type", contentType)

	// 返回文件数据
	w.Write(data)
	return nil
}

// 根据文件扩展名返回合适的 MIME 类型
func mimeByExtension(ext string) string {
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream" // 默认类型
	}
}
