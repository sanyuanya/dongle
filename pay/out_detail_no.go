package pay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sanyuanya/dongle/pay/common"
)

type OutDetailNoResponse struct {
	Mchid          string `json:"mchid"`
	OutBatchNo     string `json:"out_batch_no"`
	BatchId        string `json:"batch_id"`
	Appid          string `json:"appid"`
	OutDetailNo    string `json:"out_detail_no"`
	DetailId       string `json:"detail_id"`
	DetailStatus   string `json:"detail_status"`
	TransferAmount int    `json:"transfer_amount"`
	TransferRemark string `json:"transfer_remark"`
	FailReason     string `json:"fail_reason"`
	OpenId         string `json:"openid"`
	UserName       string `json:"user_name"`
	InitiateTime   string `json:"initiate_time"`
	UpdateTime     string `json:"update_time"`
}

func OutDetailNo(outBatchNo string, outDetailNo string) (*OutDetailNoResponse, error) {

	host := "https://api.mch.weixin.qq.com"

	path := fmt.Sprintf("/v3/transfer/batches/out-batch-no/%s/details/out-detail-no/%s", outBatchNo, outDetailNo)

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

	authorization, err := common.Signature(method, path, timestamp, nonceStr, "", privateKey)
	if err != nil {
		return nil, fmt.Errorf("无法生成签名: %v", err)
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authorization)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 打印响应详细信息
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("无法读取响应体: %v", err)
		}
		fmt.Printf("响应状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应头: %v\n", resp.Header)
		fmt.Printf("响应体: %s\n", string(respBody))
		return nil, fmt.Errorf("无法通过商家明细单号查询明细单: %v", resp.Status)
	}
	outDetailNoResponse := &OutDetailNoResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&outDetailNoResponse); err != nil {
		return nil, fmt.Errorf("无法解析响应: %v", err)
	}
	log.Printf("通过商家批次单号查询批次单: %#+v\n", outDetailNoResponse)

	return outDetailNoResponse, nil

}
