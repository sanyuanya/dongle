package tools

import (
	"fmt"
	"mime"
	"path/filepath"
)

func ValidateMimeType(fileName string) error {
	ext := filepath.Ext(fileName)
	fmt.Printf("ext: %s\n", ext)
	mimeType := mime.TypeByExtension(ext)
	fmt.Printf("mimeType: %s\n", mimeType)
	if mimeType == "" {
		return fmt.Errorf("无效的 MIME 类型: %s", ext)
	}
	return nil
}
