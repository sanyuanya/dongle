package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func Pay(c fiber.Ctx) error {

	batchesRequest := &pay.BatchesRequest{}
	err := c.Bind().Body(batchesRequest)

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	fmt.Printf("batchesRequest: %#+v\n", batchesRequest)
	resp, err := pay.Batches(batchesRequest)

	if err != nil {
		panic(fmt.Errorf("无法发起批量转账: %v", err))
	}

	return c.JSON(resp)
}
