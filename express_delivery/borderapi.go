package expressdelivery

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type KOrderApi struct {
	Kuaidicom        string `json:"kuaidicom"`
	RecManName       string `json:"recManName"`
	RecManMobile     string `json:"recManMobile"`
	RecManPrintAddr  string `json:"recManPrintAddr"`
	SendManName      string `json:"sendManName"`
	SendManMobile    string `json:"sendManMobile"`
	SendManPrintAddr string `json:"sendManPrintAddr"`
	CallBackUrl      string `json:"callBackUrl"`
	Cargo            string `json:"cargo"`
}

func BorderApi(payload *KOrderApi) error {

	url := "https://poll.kuaidi100.com/order/borderapi.do"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}

	_ = client
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("method", "bOrder")
	// mbzPBBLg6641

	return nil
}
