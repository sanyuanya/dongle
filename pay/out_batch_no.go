package pay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sanyuanya/dongle/pay/common"
)

type OutBatchNoResponse struct {
	TransferBatch      TransferBatch               `json:"transfer_batch"`
	TransferDetailList []*OutBatchNoTransferDetail `json:"transfer_detail_list,omitempty"`
}

type OutBatchNoTransferDetail struct {
	DetailId     string `json:"detail_id"`
	OutDetailNo  string `json:"out_detail_no"`
	DetailStatus string `json:"detail_status"`
}

type TransferBatch struct {
	Mchid           string `json:"mchid"`
	OutBatchNo      string `json:"out_batch_no"`
	BatchId         string `json:"batch_id"`
	Appid           string `json:"appid"`
	BatchStatus     string `json:"batch_status"`
	BatchType       string `json:"batch_type"`
	BatchName       string `json:"batch_name"`
	BatchRemark     string `json:"batch_remark"`
	CloseReason     string `json:"close_reason"`
	TotalAmount     int    `json:"total_amount"`
	TotalNum        int    `json:"total_num"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`
	SuccessAmount   int    `json:"success_amount"`
	SuccessNum      int    `json:"success_num"`
	FailAmount      int    `json:"fail_amount"`
	FailNum         int    `json:"fail_num"`
	TransferSceneId string `json:"transfer_scene_id"`
}

func OutBatchNo(outBatchNo string) (*OutBatchNoResponse, error) {

	url := "https://api.mch.weixin.qq.com/v3/transfer/batches/out-batch-no/" + outBatchNo
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
		return nil, fmt.Errorf("通过商家批次单号查询批次单: %v", resp.Status)
	}

	outBatchNoResponse := &OutBatchNoResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&outBatchNoResponse); err != nil {
		return nil, fmt.Errorf("无法解析响应: %v", err)
	}
	log.Printf("通过商家批次单号查询批次单: %#+v\n", outBatchNoResponse)

	return outBatchNoResponse, nil
}
