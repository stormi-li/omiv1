package cert

import (
	"embed"
	"os"
)

//go:embed credential/*
var certSource embed.FS

func CreatDefaultCredentialFile() {
	certFile, _ := certSource.ReadFile("credential/server.crt")

	keyFile, _ := certSource.ReadFile("credential/server.key")

	// 将嵌入的证书文件写入当前文件夹
	os.WriteFile("server.crt", certFile, 0644)

	os.WriteFile("server.key", keyFile, 0644)
}
