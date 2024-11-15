package pc

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func Pay(c fiber.Ctx) error {

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

	batchesRequest := &pay.BatchesRequest{}
	err := c.Bind().Body(batchesRequest)

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	log.Printf("batchesRequest: %#+v\n", batchesRequest)
	resp, err := pay.Batches(batchesRequest)

	if err != nil {
		panic(fmt.Errorf("无法发起批量转账: %v", err))
	}

	return c.JSON(resp)
}
