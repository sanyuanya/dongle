package common

import (
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	s, err := GenerateRandomString(32)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(s) != 32 {
		t.Fatalf("字符串长度错误")
	}
	t.Logf("success: %v", s)
}
