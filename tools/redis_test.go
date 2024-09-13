package tools

import (
	"fmt"
	"sync"
	"testing"
)

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

	wg := sync.WaitGroup{}

	// 模拟并发
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			fmt.Println("当前库存")
			res, err := rdb.DeductStock("test3", 1)
			if err != nil {
				t.Errorf("DeductStock() error = %v", err)
			}

			if !res {
				t.Errorf("DeductStock() error = %v", res)
			}
			fmt.Println("没有超卖")

		}()
	}

	wg.Wait()
}
