package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func UpdateIncome(c fiber.Ctx) error {
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

	operationId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := &entity.UpdateIncomeRequest{}

	err = c.Bind().Body(payload)

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	if payload.Shipments <= 0 {
		panic(tools.CustomError{Code: 40000, Message: "出货量不能为小于或等于0"})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	// 查询收入记录是否存在 返回收入记录
	income, err := data.GetIncomeBySnowflakeId(tx, payload.SnowflakeId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("查询收入记录失败: %v", err)})
	}

	if income == nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50004, Message: "收入记录不存在"})
	}

	if payload.Shipments == income.Shipments {
		tx.Rollback()
		panic(tools.CustomError{Code: 50004, Message: "出货量未发生变化"})
	}

	// 计算总积分
	payload.Integral = income.ProductIntegral * payload.Shipments

	err = data.UpdateIncomeBySnowflakeId(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新收入记录失败: %v", err)})
	}

	// 需要更新用户的积分 如果出货量增加则增加积分 如果出货量减少则减少积分
	if payload.Shipments > income.Shipments {
		addIntegral := payload.Integral - income.Integral
		err = data.UpdateUserIntegral(tx, income.UserId, addIntegral)
	} else {
		subIntegral := income.Integral - payload.Integral

		err = data.UpdateUserIntegral(tx, income.UserId, -subIntegral)
	}

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新用户积分失败: %v", err)})
	}

	err = data.AddOperationLog(tx, &entity.AddOperationLogRequest{
		SnowflakeId:             tools.SnowflakeUseCase.NextVal(),
		OperationId:             operationId,
		IncomeExpenseId:         income.SnowflakeId,
		UserId:                  income.UserId,
		BeforeUpdatingShipments: income.Shipments,
		AfterUpdatingShipments:  payload.Shipments,
		Summary:                 "更新客户出货量",
	})

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加操作日志失败: %v", err)})
	}

	err = tx.Commit()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("提交事务失败: %v", err)})
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "更新收入记录成功",
		Result:  struct{}{},
	})

}
