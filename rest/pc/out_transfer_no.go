package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func OutTransferNo(c fiber.Ctx) error {
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

	transferId := c.Params("transferId", "")

	if batchId == "" || transferId == "" {
		panic(tools.CustomError{Code: 40000, Message: "参数错误"})
	}

	outTransferNoResponse, err := pay.OutDetailNo(batchId, transferId)
	if err != nil {
		panic(tools.CustomError{Code: 40001, Message: fmt.Sprintf("查询失败：%v", err)})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	// 更新提现记录
	err = data.UpdateWithdrawalInfoBySnowflakeId(tx, outTransferNoResponse)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 40003, Message: fmt.Sprintf("更新失败：%v", err)})
	}

	// 如果支付失败，把提现金额退回到用户账户
	if outTransferNoResponse.DetailStatus == "FAIL" {

		withdrawal, err := data.GetWithdrawalBySnowflakeId(tx, outTransferNoResponse.OutDetailNo)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("获取提现记录失败: %v", err)
		}

		// err = data.UpdateWithdrawalStatusBySnowflakeId(tx, outTransferNoResponse.OutDetailNo, "FAIL")
		// if err != nil {
		// 	tx.Rollback()
		// 	return fmt.Errorf("更新提现状态失败: %v", err)
		// }

		err = data.AddIntegralAndWithdrawablePointsBySnowflakeId(tx, withdrawal.UserId, withdrawal.Integral)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("增加用户积分失败: %v", err)
		}

		addIncomeExpenseRequest := new(entity.AddIncomeExpenseRequest)
		addIncomeExpenseRequest.SnowflakeId = tools.SnowflakeUseCase.NextVal()
		addIncomeExpenseRequest.Summary = "退回提现积分"
		addIncomeExpenseRequest.Integral = withdrawal.Integral
		addIncomeExpenseRequest.UserId = withdrawal.UserId

		err = data.AddIncomeExpense(tx, addIncomeExpenseRequest)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("新增收支记录失败: %v", err)
		}
	}

	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})

}
