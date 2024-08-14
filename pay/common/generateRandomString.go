package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("无法生成随机字符串: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}
