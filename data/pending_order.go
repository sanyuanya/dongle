package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sanyuanya/dongle/tools"
)

func StartTicker() {
	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
		closeOrder()
		checkPendingWithdrawals()
		outTradeNo()
	}
}

func outTradeNo() {

	outTradeNoList, err := GetOrderByTradeState()

	if err != nil {
		panic(err)
	}
	for _, orderInfo := range outTradeNoList {

		resp, err := http.Get("http://localhost:3000/api/pc/outTradeNo/" + orderInfo.OutTradeNo + "/" + orderInfo.SnowflakeId)
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

		transform := &tools.Response{}
		if err := json.NewDecoder(resp.Body).Decode(&transform); err != nil {
			panic(fmt.Sprintf("获取支付状态失败: %v ", err))
		}

		if transform.Code != 0 {
			panic(fmt.Sprintf("获取支付状态失败: %v ", transform.Message))
		}

		log.Printf("获取支付状态成功: %v", transform)

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

		transform := &tools.Response{}
		if err := json.NewDecoder(resp.Body).Decode(&transform); err != nil {
			panic(fmt.Sprintf("获取支付状态失败: %v ", err))
		}

		if transform.Code != 0 {
			panic(fmt.Sprintf("获取支付状态失败: %v ", transform.Message))
		}

		log.Printf("获取支付状态成功: %v", withdrawal)
	}
}

func closeOrder() {
	cancelList, err := GetOrderExpired()
	if err != nil {
		panic(err)
	}

	for _, outTradeNo := range cancelList {
		resp, err := http.Post("http://localhost:3000/api/order/cancel/"+outTradeNo, "application/json", nil)
		if err != nil {
			panic(fmt.Sprintf("关闭订单失败: %v", err))
		}

		if resp.StatusCode != 200 {
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic(fmt.Sprintf("关闭订单失败: %v body: %#+v", err, string(body)))
			}
			panic(fmt.Sprintf("关闭订单失败: %v body: %#+v", err, string(body)))
		}

		defer resp.Body.Close()

		transform := &tools.Response{}
		if err := json.NewDecoder(resp.Body).Decode(&transform); err != nil {
			panic(fmt.Sprintf("关闭订单失败：%v ", err))
		}

		if transform.Code != 0 {
			panic(fmt.Sprintf("关闭订单失败：%v ", transform.Message))
		}

		log.Printf("关闭订单成功：%v", outTradeNo)
	}
}
