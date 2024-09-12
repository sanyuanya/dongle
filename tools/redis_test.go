package tools

import "testing"

func TestNewRedisClient(t *testing.T) {
	rdb := Redis{}

	err := rdb.NewClient()
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
	}

	err = rdb.SetSKUStock("test3", 998)

	if err != nil {
		t.Errorf("SetSKUStock() error = %v", err)
	}

}
