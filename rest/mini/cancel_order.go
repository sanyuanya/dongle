package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func CancelOrder(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}
	outTradeNo := c.Params("orderId")
	err = pay.CloseOrder(outTradeNo)
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("关闭订单失败: %v", err)})
	}

	// 更新订单状态

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	updateOrderByOutTradeNo := &entity.UpdateOrderByOutTradeNo{
		Status:     99,
		OutTradeNo: outTradeNo,
	}

	err = data.UpdateOrderByOutTradeNo(tx, updateOrderByOutTradeNo)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新订单状态失败: %v", err)})
	}
	return c.JSON(tools.Response{
		Code:    0,
		Message: "关闭订单成功",
		Result:  struct{}{},
	})

}
