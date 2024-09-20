package data

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func StartTicker() {
	ticker := time.NewTicker(3 * time.Minute)

	for range ticker.C {
		checkPendingWithdrawals()
		outTradeNo()
	}
}

func outTradeNo() {

	outTradeNoList, err := GetOrderByTradeState()

	if err != nil {
		panic(err)
	}
	for _, outTradeNo := range outTradeNoList {

		resp, err := http.Get("http://localhost:3000/api/pc/outTradeNo/" + outTradeNo)
		if err != nil {
			panic(fmt.Sprintf("获取支付状态失败: %v", err))
		}

		if resp.StatusCode != 200 {
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic(fmt.Sprintf("获取支付状态失败: %v body: %#+v", err, string(body)))
			}
			panic(fmt.Sprintf("获取支付状态失败: %v body: %#+v", err, string(body)))
		}

		defer resp.Body.Close()
	}
}

func checkPendingWithdrawals() {

	withdrawalList, err := GetWithdrawalByPaymentStatusIsFailAndPaymentStatusIsSuccess()

	if err != nil {
		panic(err)
	}

	for _, withdrawal := range withdrawalList {

		resp, err := http.Get("http://localhost:3000/api/pc/batch/" + withdrawal.PayId + "/transfer/" + withdrawal.SnowflakeId)
		if err != nil {
			panic(fmt.Sprintf("获取支付状态失败: %v", err))
		}

		if resp.StatusCode != 200 {
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic(fmt.Sprintf("获取支付状态失败: %v body: %#+v", err, string(body)))
			}
			panic(fmt.Sprintf("获取支付状态失败: %v body: %#+v", err, string(body)))
		}

		defer resp.Body.Close()
	}

}
