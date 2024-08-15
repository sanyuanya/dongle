package pay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

	host := "https://api.mch.weixin.qq.com"

	path := "/v3/transfer/batches/out-batch-no"
	u, err := url.Parse(host)

	if err != nil {
		return nil, fmt.Errorf("无法解析 URL: %v", err)
	}

	u.Path, err = url.JoinPath(u.Path, path, outBatchNo)
	if err != nil {
		return nil, fmt.Errorf("无法拼接 URL: %v", err)
	}

	query := url.Values{}

	query.Add("need_query_detail", "true")
	query.Add("offset", "0")
	query.Add("limit", "100")
	query.Add("detail_status", "ALL")

	u.RawQuery = query.Encode()

	method := http.MethodGet

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	nonceStr, err := common.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("无法生成随机字符串: %v", err)
	}

	env := os.Getenv("ENVIRONMENT")

	certPath := ""
	switch env {
	case "production":
		certPath = "/cert"
	default:
		certPath = "/Users/sanyuanya/hjworkspace/go_dev/dongle_new/pay/cert"
	}

	privateFilePath := fmt.Sprintf("%s/apiclient_key.pem", certPath)
	privateKey, err := common.ReadPrivateKey(privateFilePath)
	if err != nil {
		return nil, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	signUrl := &url.URL{Path: fmt.Sprintf("%s/%s", path, outBatchNo), RawQuery: u.RawQuery}

	authorization, err := common.Signature(method, signUrl.String(), timestamp, nonceStr, "", privateKey)
	if err != nil {
		return nil, fmt.Errorf("无法生成签名: %v", err)
	}

	fmt.Println(u.String())
	req, err := http.NewRequest(method, u.String(), nil)

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
		return nil, fmt.Errorf("无法通过商家批次单号查询批次单: %v", resp.Status)
	}

	outBatchNoResponse := &OutBatchNoResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&outBatchNoResponse); err != nil {
		return nil, fmt.Errorf("无法解析响应: %v", err)
	}
	log.Printf("通过商家批次单号查询批次单: %#+v\n", outBatchNoResponse)

	return outBatchNoResponse, nil
}
