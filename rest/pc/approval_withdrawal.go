package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func ApprovalWithdrawal(c fiber.Ctx) error {

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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	approvalWithdrawalRequest := &entity.ApprovalWithdrawalRequest{}

	err = c.Bind().Body(approvalWithdrawalRequest)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	err = data.ApprovalWithdrawal(approvalWithdrawalRequest)
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法审批提现: %v", err)})
	}

	
	return c.JSON(tools.Response{
		Code:    0,
		Message: "审批提现成功",
		Result:  struct{}{},
	})
}
