package pay

import (
	"fmt"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	host := "https://api.mch.weixin.qq.com"

	outBatchNo := "123"
	path := "/v3/transfer/batches/out-batch-no/" + outBatchNo

	u, err := url.Parse(host)

	if err != nil {
		t.Errorf("无法解析 URL: %v", err)
	}

	u.Path, err = url.JoinPath(u.Path, path, outBatchNo)
	if err != nil {
		t.Errorf("无法拼接 URL: %v", err)
	}

	query := url.Values{}

	query.Add("need_query_detail", "false")
	query.Add("offset", "0")
	query.Add("limit", "100")
	query.Add("detail_status", "ALL")

	u.RawQuery = query.Encode()

	s := &url.URL{Path: path, RawQuery: u.RawQuery}
	fmt.Printf("url: %s", s.String())
}
