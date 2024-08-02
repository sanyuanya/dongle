package rest

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
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	approvalWithdrawalRequest := &entity.ApprovalWithdrawalRequest{}

	err = c.Bind().Body(approvalWithdrawalRequest)
	if err != nil {
		panic(fmt.Errorf("无法绑定请求体: %v", err))
	}

	err = data.ApprovalWithdrawal(approvalWithdrawalRequest)
	if err != nil {
		panic(fmt.Errorf("无法审批提现: %v", err))
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "审批提现成功",
		Result:  struct{}{},
	})
}
