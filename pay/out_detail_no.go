package pay

import (
	"encoding/json"
	"fmt"
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
	url := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/transfer/batches/out-batch-no/%s/details/out-detail-no/%s", outBatchNo, outDetailNo)
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

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("通过商家明细单号查询明细单: %v", resp.Status)
	}
	outDetailNoResponse := &OutDetailNoResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&outDetailNoResponse); err != nil {
		return nil, fmt.Errorf("无法解析响应: %v", err)
	}
	log.Printf("通过商家批次单号查询批次单: %#+v\n", outDetailNoResponse)

	return outDetailNoResponse, nil

}
