package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/entity"
	expressdelivery "github.com/sanyuanya/dongle/express_delivery"
	"github.com/sanyuanya/dongle/tools"
)

func Label(c fiber.Ctx) error {
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

	payload := &entity.LabelOrderRequest{}
	if err := c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	resp, err := expressdelivery.LabelOrder(payload)

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("下单失败: %v", err)})
	}

	return c.JSON(resp)
}
