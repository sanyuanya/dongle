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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	applyForWithdrawal := new(entity.ApplyForWithdrawalRequest)
	err = c.Bind().Body(applyForWithdrawal)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	// 判断是否是白名单用户
	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	err = data.IsWhite(tx, snowflakeId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", err)})
	}

	// 判断用户当前所属积分是否大于等于提现积分
	err = data.IsIntegralWithdraw(tx, snowflakeId, applyForWithdrawal.Integral)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", err)})
	}

	// 判断 可提现积分是否大于等于提现积分
	userDetail, err := data.GetUserDetailBySnowflakeID(tx, snowflakeId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", err)})
	}

	if userDetail.WithdrawablePoints < applyForWithdrawal.Integral {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", "当前提现积分大于可提现积分")})
	}

	if applyForWithdrawal.Integral <= 0 {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", "提现积分必须大于0")})
	}

	applyForWithdrawal.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	applyForWithdrawal.UserId = snowflakeId

	err = data.ApplyForWithdrawal(tx, applyForWithdrawal)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", err)})
	}

	// 扣除用户积分和可提现积分
	err = data.DeductUserIntegralAndWithdrawablePointsBySnowflakeId(tx, snowflakeId, applyForWithdrawal.Integral)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", err)})
	}

	if applyForWithdrawal.WithdrawalMethod == "alipay" {
		err = data.UpdateUserAlipayAccountBySnowflakeId(
			tx,
			applyForWithdrawal.UserId,
			applyForWithdrawal.AlipayAccount,
		)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法申请提现: %v", err)})
		}
	}
	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "申请提现成功",
		Result:  struct{}{},
	})
}
