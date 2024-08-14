package pay

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/pay/common"
)

type BatchesHeader struct {
	Authorization   string `json:"Authorization"`
	Accept          string `json:"Accept"`
	ContentType     string `json:"Content-Type"`
	WechatpaySerial string `json:"Wechatpay-Serial"`
}

type BatchesRequest struct {
	AppId              string           `json:"appid"`
	OutBatchNo         string           `json:"out_batch_no"`
	BatchName          string           `json:"batch_name"`
	BatchRemark        string           `json:"batch_remark"`
	TotalAmount        int              `json:"total_amount"`
	TotalNum           int              `json:"total_num"`
	TransferDetailList []TransferDetail `json:"transfer_detail_list"`
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

func Batches() error {

	url := "https://api.mch.weixin.qq.com/v3/transfer/batches"

	method := "POST"

	timestamp := time.Now().Unix()

	nonceStr := "5K8264ILTKCH16CQ2502SI8ZNMTM67VS"

	body := &BatchesRequest{
		AppId:       "wx8888888888888888",
		OutBatchNo:  "plfk2020042013",
		BatchName:   "2019年1月深圳分部报销单",
		BatchRemark: "2019年1月深圳分部报销单",
		TotalAmount: 4000,
		TotalNum:    2,
		TransferDetailList: []TransferDetail{
			{
				OutDetailNo:    "x23zy545Bd5436",
				TransferAmount: 4000,
				TransferRemark: "深圳分部报销",
				OpenId:         "o-MYE5dGdI3cFz2t7zjDzjDx5K8",
				UserName:       "张三",
			},
			{
				OutDetailNo:    "x23zy545Bd5437",
				TransferAmount: 4000,
				TransferRemark: "深圳分部报销",
				OpenId:         "o-MYE5dGdI3cFz2t7zjDzjDx5K8",
				UserName:       "李四",
			},
		},
	}

	payloadByte, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("无法序列化请求体: %v", err)
	}

	privateKey, err := common.ReadPrivateKey("apiclient_key.pem")
	if err != nil {
		return fmt.Errorf("无法读取私钥文件: %v", err)
	}

	authorization, err := common.Signature(method, url, fmt.Sprintf("%d", timestamp), nonceStr, string(payloadByte), privateKey)
	if err != nil {
		return fmt.Errorf("无法生成签名: %v", err)
	}

	serialNo := "17BDDF6F46451DE2C953B628B76D4458B00CF054"

	publicKey, err := common.ReadPublicKey("apiclient_cert.pem")
	if err != nil {
		return fmt.Errorf("无法读取公钥文件: %v", err)
	}

	encrypt, err := common.Encrypt([]byte(serialNo), publicKey)
	if err != nil {
		return fmt.Errorf("无法加密签名: %v", err)
	}

	_ = &BatchesHeader{
		Accept:          "application/json",
		ContentType:     "application/json",
		Authorization:   authorization,
		WechatpaySerial: string(encrypt),
	}

}
