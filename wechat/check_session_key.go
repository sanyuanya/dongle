package wechat

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
)

func generateHMACSHA256(sessionKey string, message string) string {
	// 创建一个新的 HMAC 使用 SHA256 哈希算法
	h := hmac.New(sha256.New, []byte(sessionKey))

	// 写入消息数据
	h.Write([]byte(message))

	// 计算 HMAC
	return hex.EncodeToString(h.Sum(nil))
}

func CheckSessionKey(openid, accessToken, sessionKey string) (error, error) {

	baseURL := "https: //api.weixin.qq.com/wxa/checksession"

	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("openid", openid)
	q.Set("access_token", accessToken)
	q.Set("signature", generateHMACSHA256(sessionKey, ""))
	q.Set("sig_method", "hmac_sha256")

	u.RawQuery = q.Encode()

	// 创建一个自定义的 http.Client，跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(u.String())

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	code2SessionResp := Code2SessionResp{}
	if err := json.NewDecoder(resp.Body).Decode(&code2SessionResp); err != nil {
		return nil, err
	}
	return nil, nil
}
