package tools

import "testing"

func TestNewRedisClient(t *testing.T) {
	rdb := Redis{}

	err := rdb.NewClient()
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
	}

}
