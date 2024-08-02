package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/snowflake"
)

func ApplyForWithdrawal(c fiber.Ctx) error {
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

	if role != "user" {
		panic(fmt.Errorf("未经授权"))
	}

	snowflakeId, err := strconv.ParseInt(snowflakeIdStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err))
	}

	applyForWithdrawal := new(entity.ApplyForWithdrawalRequest)
	err = c.Bind().Body(applyForWithdrawal)
	if err != nil {
		panic(fmt.Errorf("无法绑定请求体: %v", err))
	}

	err = data.IsWhite(snowflakeId)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	isWithdraw, err := data.IsIntegralWithdraw(snowflakeId, applyForWithdrawal.Integral)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	if !isWithdraw {
		panic(fmt.Errorf("无法申请提现：%v", err))
	}

	applyForWithdrawal.SnowflakeId = snowflake.SnowflakeUseCase.NextVal()
	applyForWithdrawal.UserId = snowflakeId

	err = data.ApplyForWithdrawal(applyForWithdrawal)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	if applyForWithdrawal.WithdrawalMethod == "alipay" {
		err = data.UpdateUserAlipayAccountBySnowflakeID(applyForWithdrawal.UserId, applyForWithdrawal.AlipayAccount)
		if err != nil {
			panic(fmt.Errorf("无法申请提现: %v", err))
		}
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "申请提现成功",
		Result:  struct{}{},
	})
}
