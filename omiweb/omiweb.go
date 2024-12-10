package web

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var Open = false

type Web struct {
	EmbeddedSource *embed.FS
}

const IndexPath = "static/templates/index.html"

func NewWebServer(embeddedSource *embed.FS) *Web {
	if embeddedSource != nil {
		_, err := embeddedSource.Open(IndexPath)
		if err != nil {
			panic(err)
		}
	}
	Open = true
	return &Web{
		EmbeddedSource: embeddedSource,
	}
}

const SourcePath = "static"

func (web *Web) GenerateTemplate() {
	copyEmbeddedFiles(SourcePath)
}

const favicon = "/favicon.ico"

func removeBeforeStatic(input string) string {
	index := strings.Index(input, SourcePath)
	if index == -1 {
		// 如果没有找到 "/static"，返回原始字符串
		return input
	}
	// 返回从 "/static" 开始的子字符串
	if index > 1 && input[index-1] != '.' {
		return input
	}
	return input[index:]
}

const TemplatesPath = "static/templates"

func (web *Web) ServeWeb(w http.ResponseWriter, r *http.Request) bool {
	filePath := r.URL.Path
	if strings.ToLower(filepath.Ext(filePath)) == ".html" {
		filePath = TemplatesPath + filePath
	}
	filePath = removeBeforeStatic(filePath)

	if filePath == favicon {
		filePath = SourcePath + filePath
	}

	if filePath == "/" {
		filePath = IndexPath
	}

	var data []byte
	var err error
	if web.EmbeddedSource != nil {
		data, err = web.EmbeddedSource.ReadFile(filePath)
	} else {
		data, err = os.ReadFile(filePath)
	}
	// 检查是否读取失败
	if err != nil {
		return false
	}

	WriterHeader(w, r)

	// 返回文件数据
	w.Write(data)
	return true
}

func WriterHeader(w http.ResponseWriter, r *http.Request) {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(r.URL.Path))
	// 设置Content-Type
	contentType := mimeByExtension(ext)
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
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
		return ""
	}
}
