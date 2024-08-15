package pay

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sanyuanya/dongle/pay/common"
)

func Certificates() (error, error) {
	host := "https://api.mch.weixin.qq.com"

	path := "/v3/certificates"

	url := fmt.Sprintf("%s%s", host, path)

	method := http.MethodGet

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	nonceStr, err := common.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("无法生成随机字符串: %v", err)
	}

	privateKey, err := common.ReadPrivateKey("apiclient_key.pem")
	if err != nil {
		return nil, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	authorization, err := common.Signature(method, url, timestamp, nonceStr, "", privateKey)
	if err != nil {
		return nil, fmt.Errorf("无法生成签名: %v", err)
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %v", err)
	}

	defer resp.Body.Close()

	// 打印响应详细信息
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("无法读取响应体: %v", err)
	}
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应头: %v\n", resp.Header)
	fmt.Printf("响应体: %s\n", string(respBody))

	return nil, nil
}
