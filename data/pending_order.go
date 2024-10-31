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
	log.Println("开启定时器 读取交易状态订单 ................")
	outTradeNoList, err := GetOrderByTradeState()
	if err != nil {
		panic(fmt.Sprintf("从数据库获取交易状态订单时出错 %v", err))
	}
	log.Printf("从数据库获取交易状态订单数据 %v \n", outTradeNoList)
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

	log.Println("开启定时器 读取申请提现状态订单 ................")

	withdrawalList, err := GetWithdrawalByPaymentStatusIsFailAndPaymentStatusIsSuccess()

	if err != nil {
		panic(fmt.Sprintf("从数据库查询申请提现订单时出错 %v", err))
	}

	for _, withdrawal := range withdrawalList {

		resp, err := http.Get("http://localhost:3000/api/pc/batch/" + withdrawal.PayId + "/transfer/" + withdrawal.SnowflakeId)
		if err != nil {
			panic(fmt.Sprintf("获取申请提现支付状态失败: %v", err))
		}

		if resp.StatusCode != 200 {
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				panic(fmt.Sprintf("获取申请提现支付状态失败: %v body: %#+v", err, string(body)))
			}
			panic(fmt.Sprintf("获取申请提现支付状态失败: %v body: %#+v", err, string(body)))
		}

		defer resp.Body.Close()

		transform := &tools.Response{}
		if err := json.NewDecoder(resp.Body).Decode(&transform); err != nil {
			panic(fmt.Sprintf("获取申请提现支付状态失败: %v ", err))
		}

		if transform.Code != 0 {
			panic(fmt.Sprintf("获取申请提现支付状态失败: %v ", transform.Message))
		}

		log.Printf("获取申请提现支付状态成功: %v", withdrawal)
	}
}

func closeOrder() {
	log.Println("开启定时器 读取待支付状态订单 ................")

	cancelList, err := GetOrderExpired()
	if err != nil {
		panic(fmt.Sprintf("读取超时未支付订单出错 %v", err))
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
