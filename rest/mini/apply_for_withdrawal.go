package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func ApplyForWithdrawal(c fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(tools.Response{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	applyForWithdrawal := new(entity.ApplyForWithdrawalRequest)
	err = c.Bind().Body(applyForWithdrawal)
	if err != nil {
		panic(fmt.Errorf("无法绑定请求体: %v", err))
	}

	// 判断是否是白名单用户
	err = data.IsWhite(snowflakeId)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	// 判断用户当前所属积分是否大于等于提现积分
	err = data.IsIntegralWithdraw(snowflakeId, applyForWithdrawal.Integral)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	// 判断 可提现积分是否大于等于提现积分
	userDetail, err := data.GetUserDetailBySnowflakeID(snowflakeId)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	if userDetail.WithdrawablePoints < applyForWithdrawal.Integral {
		panic(fmt.Errorf("无法申请提现: %v", "当前提现积分大于可提现积分"))
	}

	applyForWithdrawal.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	applyForWithdrawal.UserId = snowflakeId

	err = data.ApplyForWithdrawal(applyForWithdrawal)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	// 扣除用户积分和可提现积分
	err = data.DeductUserIntegralAndWithdrawablePointsBySnowflakeId(snowflakeId, applyForWithdrawal.Integral)

	if err != nil {
		panic(fmt.Errorf("无法申请提现: %v", err))
	}

	if applyForWithdrawal.WithdrawalMethod == "alipay" {
		err = data.UpdateUserAlipayAccountBySnowflakeId(
			applyForWithdrawal.UserId,
			applyForWithdrawal.AlipayAccount,
		)
		if err != nil {
			panic(fmt.Errorf("无法申请提现: %v", err))
		}
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "申请提现成功",
		Result:  struct{}{},
	})
}
