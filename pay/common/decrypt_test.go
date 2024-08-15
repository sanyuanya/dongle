package common

import (
	"fmt"
	"testing"
)

func TestDecrypt(t *testing.T) {

	str := "cRwKXziuRM6uCEVLLp7YRJWvU6256+c+cn4kcwZI688rE27adpaCKT5HHIeoXsOgGWVaGcUpgF6QDa9eEInZr2zdA7lyMZL4yUezucVjkUDoFEH8NPzIDwDNLuimbqWR2CXmdRi+b69XROc0n+lF2iQw8OruqhFkrAgbVINZHF5DN6OOQeUpvK6cyqzkjdhQef+HH3P6b5L3wYha8HTfvzhE9VowwNkqbfQk6XlOZvEfplbHejxHYx9psOtte+UR/lqPhafV8I0nz1Sw8P4JJb9CI673oNKGsoQ761efgksrPKWwDiXTLjvvdt8SuOdw2s3J7wKCKt4uNnyzq2VNgQ=="
	e, err := Decrypt(str)

	if err != nil {
		t.Errorf("err: %v", err)
	}
	fmt.Println(e)
}
