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

	"github.com/sanyuanya/dongle/entity"
)

func PollQuery(param *entity.PollQueryRequest) (*entity.PollQueryResponse, error) {
	key := "mbzPBBLg6641"
	customer := "8A4E62031DA6A56D23825E817128652F"
	apiURL := "https://poll.kuaidi100.com/poll/query.do"
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf("快递100 实时快递查询下单接口 json marshal PollQueryRequest struct error %v", err)
	}
	tempSign := string(paramBytes) + key + customer
	hash := md5.New()
	hash.Write([]byte(tempSign))
	sign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	data := url.Values{}
	data.Set("customer", customer)
	data.Set("param", string(paramBytes))
	data.Set("sign", sign)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("快递100 实时快递查询下单接口http发起请求出错:%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("快递100 实时快递查询下单接口请求失败:%s", resp.Status)
	}
	var pollQueryResponse *entity.PollQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&pollQueryResponse); err != nil {
		return nil, fmt.Errorf("快递100 实时快递查询下单接口解析响应失败:%s", err)
	}
	return pollQueryResponse, nil
}
