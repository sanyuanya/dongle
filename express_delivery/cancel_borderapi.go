package expressdelivery

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func CancelBorderApi(payload *entity.CancelKOrderApiRequest) (*entity.CancelKOrderApiResponse, error) {

	apiURL := "https://poll.kuaidi100.com/order/borderapi.do"

	key := "mbzPBBLg6641"
	secret := "f969f49a93dc45979478aece402b0264"
	t := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("快递100 商家取消寄件请求下单接口 json marshal KOrderApiRequestParam struct error : %v", err)
	}

	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s%s%s%s", string(payloadBytes), t, key, secret)))
	md5String := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))

	data := url.Values{}

	data.Set("key", key)
	data.Set("method", "bOrder")
	data.Set("t", t)
	data.Set("param", string(payloadBytes))
	data.Set("sign", md5String)

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("快递100 商家取消寄件请求下单接口 http 发起请求 出错: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("快递100 商家取消寄件请求下单接口 请求失败: %s", resp.Status)
	}

	var cancelKOrderApiResponse *entity.CancelKOrderApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&cancelKOrderApiResponse); err != nil {
		return nil, fmt.Errorf("快递100 商家取消寄件请求下单接口 解析响应失败：%v", err)
	}

	return cancelKOrderApiResponse, nil
}
