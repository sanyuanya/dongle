package pc

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	expressdelivery "github.com/sanyuanya/dongle/express_delivery"
	"github.com/sanyuanya/dongle/tools"
)

func Shipping(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			var code int
			var message string
			switch e := err.(type) {
			case tools.CustomError:
				code = e.Code
				message = e.Message
			case error:
				code = 50001
				message = e.Error()
			default:
				code = 50002
				message = fmt.Sprintf("%v", e)
			}
			c.JSON(tools.Response{
				Code:    code,
				Message: message,
				Result:  struct{}{},
			})
		}
	}()

	shipping := new(entity.Shipping)

	if err := c.Bind().Body(shipping); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	for _, v := range shipping.OrderList {
		fmt.Printf("shipping.OrderList: %v\n", v)

		orderCommodityList, err := data.GetOrderCommodityList(tx, v)
		if err != nil {
			log.Printf("查询订单所购商品失败: %v", err)
			panic(tools.CustomError{Code: 50007, Message: "查询订单所购商品失败,请联系管理员"})
		}

		if len(orderCommodityList) == 0 {
			panic(tools.CustomError{Code: 50007, Message: "查询订单失败"})
		}

		border := &entity.KOrderApiRequestParam{
			Kuaidicom:        "jd",
			RecManName:       orderCommodityList[0].Consignee,
			RecManMobile:     orderCommodityList[0].PhoneNumber,
			RecManPrintAddr:  fmt.Sprintf("%s%s", orderCommodityList[0].Location, orderCommodityList[0].DetailedAddress),
			SendManName:      "赵世龙",
			SendManMobile:    "15135002301",
			SendManPrintAddr: "山西省太原市小店区长治306号1幢B座11层1116室",
			Cargo:            "电子产品",
			Remark:           orderCommodityList[0].OrderId,
			CallBackUrl:      "https://www.iotpeachcloud.com/api/order/orderCallback",
		}

		resp, err := expressdelivery.BorderApi(border)
		if err != nil {
			log.Printf("商家寄件下单失败 下单参数:%#+v,  失败原因: %v", v, err)
			panic(tools.CustomError{Code: 50007, Message: "商家寄件下单失败"})
		}

		if !resp.Result {
			log.Printf("商家寄件下单失败返回编码%s返回报文描述%s", resp.ReturnCode, resp.Message)
			panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("商家寄件下单失败描述%s", resp.Message)})
		}

		addShippingRequest := entity.AddShippingRequest{
			SnowflakeId:  tools.SnowflakeUseCase.NextVal(),
			TaskId:       resp.Data.TaskId,
			OrderId:      border.Remark,
			ThirdOrderId: resp.Data.OrderId,
			OrderNumber:  resp.Data.Kuaidinum,
		}

		if err := data.AddShipping(tx, &addShippingRequest); err != nil {
			tx.Rollback()
			log.Printf("保存商家寄件信息失败 请求参数 %#+v, 失败原因： %v", addShippingRequest, err)
			continue
		}
		updateOrderStatusRequest := &entity.UpdateOrderStatusRequest{
			OrderId: addShippingRequest.OrderId,
			Status:  3,
		}

		if err := data.UpdateOrderStatus(tx, updateOrderStatusRequest); err != nil {
			tx.Rollback()
			log.Printf("更新订单状态失败 失败原因: %v", err)
		}
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "发货成功",
		Result:  map[string]any{},
	})
}
