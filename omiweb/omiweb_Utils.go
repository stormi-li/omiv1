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
	destDir := dir
	// 遍历嵌入文件系统中的所有文件
	err := fs.WalkDir(srcFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录，仅复制文件
		if d.IsDir() {
			return nil
		}

		// 确保目标文件夹路径
		destPath := filepath.Join(destDir, filepath.Base(path))
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		// 检查目标文件是否已存在
		if _, err := os.Stat(destPath); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}

		// 读取嵌入文件内容
		content, err := srcFS.ReadFile(path)
		if err != nil {
			return err
		}

		// 写入文件到目标文件夹
		if err := os.WriteFile(destPath, content, os.ModePerm); err != nil {
			return err
		}

		return nil
	})

	return err
}
