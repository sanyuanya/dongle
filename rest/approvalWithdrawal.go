package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
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

	authorization := c.Get("Authorization")

	token, err := ValidateToken(authorization)

	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	claims := token.Claims.(jwt.MapClaims)

	snowflakeIdStr, ok := claims["snowflake_id"].(string)
	if !ok {
		panic(fmt.Errorf("snowflake_id is not a string"))
	}

	role, ok := claims["role"].(string)
	if !ok {
		panic(fmt.Errorf("role is not a string"))
	}

	if role != "admin" {
		panic(fmt.Errorf("未经授权"))
	}

	snowflakeId, err := strconv.ParseInt(snowflakeIdStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err))
	}
	_ = snowflakeId

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
