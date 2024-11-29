package web

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed templateSource/static/*
var templateSource embed.FS

func copyEmbeddedFiles(dir string) error {
	srcFS := templateSource
	srcPath := "templateSource/static" // 嵌入的起始目录
	destDir := dir

	// 遍历嵌入文件系统中的所有文件和目录
	err := fs.WalkDir(srcFS, srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 构造目标路径，去除嵌入路径前缀
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		// 如果是目录，创建目录
		if d.IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return err
			}
			return nil
		}

		// 如果是文件，检查并复制文件
		// 检查目标文件是否已存在
		if _, err := os.Stat(destPath); err == nil {
			// 文件已存在，跳过
			return nil
		} else if !os.IsNotExist(err) {
			// 其他错误
			return err
		}

		// 读取嵌入文件内容
		content, err := srcFS.ReadFile(path)
		if err != nil {
			return err
		}

		// 写入文件到目标路径
		if err := os.WriteFile(destPath, content, os.ModePerm); err != nil {
			return err
		}

		return nil
	})

	return err
}
