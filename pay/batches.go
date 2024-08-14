package pay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sanyuanya/dongle/pay/common"
)

type BatchesRequest struct {
	AppId              string            `json:"appid"`
	OutBatchNo         string            `json:"out_batch_no"`
	BatchName          string            `json:"batch_name"`
	BatchRemark        string            `json:"batch_remark"`
	TotalAmount        int               `json:"total_amount"`
	TotalNum           int               `json:"total_num"`
	TransferDetailList []*TransferDetail `json:"transfer_detail_list"`
}

type TransferDetail struct {
	OutDetailNo    string `json:"out_detail_no"`
	TransferAmount int    `json:"transfer_amount"`
	TransferRemark string `json:"transfer_remark"`
	OpenId         string `json:"openid"`
	UserName       string `json:"user_name"`
}

type BatchesResponse struct {
	OutBatchNo  string `json:"out_batch_no"`
	BatchId     string `json:"batch_id"`
	CreateTime  string `json:"create_time"`
	BatchStatus string `json:"batch_status"`
}

func Batches(body *BatchesRequest) (*BatchesResponse, error) {

	url := "https://api.mch.weixin.qq.com/v3/transfer/batches"

	method := http.MethodPost

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	nonceStr, err := common.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("无法生成随机字符串: %v", err)
	}

	payloadByte, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("无法序列化请求体: %v", err)
	}

	privateKey, err := common.ReadPrivateKey("apiclient_key.pem")
	if err != nil {
		return nil, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	authorization, err := common.Signature(method, url, timestamp, nonceStr, string(payloadByte), privateKey)
	if err != nil {
		return nil, fmt.Errorf("无法生成签名: %v", err)
	}

	serialNo := "17BDDF6F46451DE2C953B628B76D4458B00CF054"

	publicKey, err := common.ReadPublicKey("apiclient_cert.pem")
	if err != nil {
		return nil, fmt.Errorf("无法读取公钥文件: %v", err)
	}

	encrypt, err := common.Encrypt([]byte(serialNo), publicKey)
	if err != nil {
		return nil, fmt.Errorf("无法加密签名: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadByte))

	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Wechatpay-Serial", string(encrypt))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("无法发起商家转账: %v", resp.Status)
	}
	batchesResponse := &BatchesResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&batchesResponse); err != nil {
		return nil, fmt.Errorf("无法解析响应: %v", err)
	}
	log.Printf("发起商家转账响应: %#+v\n", batchesResponse)

	return batchesResponse, nil

}
