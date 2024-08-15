package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func OutBatchNo(c fiber.Ctx) error {

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

	batchId := c.Params("batchId", "")

	if batchId == "" {
		panic(tools.CustomError{Code: 40000, Message: "参数错误"})
	}

	outBatchNoResponse, err := pay.OutBatchNo(batchId)
	if err != nil {
		panic(tools.CustomError{Code: 40001, Message: fmt.Sprintf("查询失败：%v", err)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	err = data.UpdatePay(tx, outBatchNoResponse.TransferBatch)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 40002, Message: fmt.Sprintf("更新失败：%v", err)})
	}

	data.Commit(tx)
	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})
}
