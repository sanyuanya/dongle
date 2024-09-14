package pay

import "testing"

func TestPay(t *testing.T) {
	pay := JsApiRequest{
		AppId:       "wx370126c8bcf8d00c",
		Mchid:       "1682195529",
		Description: "image形象店-深圳腾大-QQ公仔",
		OutTradeNo:  "1217752501201407033233368018",
		Attach:      "自定义数据",
		Amount: Amount{
			Total:    1,
			Currency: "CNY",
		},
		NotifyUrl: "https://www.weixin.qq.com/wxpay/pay.php",
		Payer: Payer{
			OpenId: "ozxx67aQHWju7uONzzz0kaiSIHxw",
		},
	}

	jsApiResponse, err := JsApi(&pay)
	if err != nil {
		t.Errorf("TestPay failed, err: %v", err)
	}

	t.Logf("jsApiResponse: %v", jsApiResponse)

}
