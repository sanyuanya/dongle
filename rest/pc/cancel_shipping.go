package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	expressdelivery "github.com/sanyuanya/dongle/express_delivery"
	"github.com/sanyuanya/dongle/tools"
)

func CancelShipping(c fiber.Ctx) error {

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

	if _, err := tools.ValidateUserToken(c.Get("Authorization"), "user"); err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}
	orderId := c.Params("orderId", "")

	request := new(entity.CancelBorderApiRequest)

	if err := c.Bind().Body(request); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	if len(request.CancelMsg) > 30 {
		panic(tools.CustomError{Code: 40000, Message: "取消原因-长度超出限制"})
	}

	if request.CancelMsg == "" {
		request.CancelMsg = "暂时不寄件了!!!"
	}
	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	shipping, err := data.GetShippingByOrderId(tx, orderId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("暂未获取到物流信息,请稍后再试: %v", err)})
	}

	if shipping == nil {
		return c.JSON(tools.Response{
			Code:    0,
			Message: "暂未获取到物流信息,请稍后再试",
			Result:  struct{}{},
		})
	}

	payload := &entity.CancelKOrderApiRequest{
		TaskId:    shipping.TaskId,
		OrderId:   shipping.ThirdOrderId,
		CancelMsg: request.CancelMsg,
	}
	resp, err := expressdelivery.CancelBorderApi(payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("取消寄件失败: %v", err)})
	}

	if resp.Result {
		if err := data.UpdateOrderStatus(tx, &entity.UpdateOrderStatusRequest{OrderId: shipping.OrderId, Status: 100}); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("取消寄件-更新订单状态失败: %v", err)})
		}
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: resp.Message,
		Result:  struct{}{},
	})
}
