package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func OutTradeNo(c fiber.Ctx) error {
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

	outTradeNo := c.Params("outTradeNo", "")

	if outTradeNo == "" {
		panic(tools.CustomError{Code: 40000, Message: "参数错误"})
	}

	outTradeNoResponse, err := pay.OutTradeNo(outTradeNo)
	if err != nil {
		panic(tools.CustomError{Code: 40001, Message: fmt.Sprintf("查询失败：%v", err)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	// 更新订单
	err = data.UpdateOrder(tx, outTradeNoResponse)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新订单失败: %v", err)})
	}

	err = tx.Commit()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("提交事务失败: %v", err)})
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result:  outTradeNoResponse,
	})

}
