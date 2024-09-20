package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/tools"
)

func GetOrderDetail(c fiber.Ctx) error {
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

	orderId := c.Params("orderId")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	order, err := data.GetOrderDetail(tx, orderId)

	if err != nil {
		panic(tools.CustomError{Code: 50004, Message: fmt.Sprintf("无法获取订单详情: %v", err)})
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取订单详情成功",
		Result:  order,
	})
}
