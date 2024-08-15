package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/pay"
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

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	for _, snowflakeId := range approvalWithdrawalRequest.ApprovalList {

		err = data.ApprovalWithdrawal(tx, snowflakeId, approvalWithdrawalRequest.Rejection, approvalWithdrawalRequest.LifeCycle)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法审批提现: %v", err)})
		}

		if approvalWithdrawalRequest.LifeCycle == 2 {
			withdrawal, err := data.GetWithdrawalBySnowflakeId(tx, snowflakeId)
			if err != nil {
				data.Rollback(tx)
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取提现记录失败: %v", err)})
			}
			err = data.AddIntegralAndWithdrawablePointsBySnowflakeId(tx, withdrawal.UserId, withdrawal.Integral)
			if err != nil {
				data.Rollback(tx)
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("增加用户积分失败: %v", err)})
			}
		} else {
			transferDetailList := []*pay.TransferDetail{}
			transferDetail, err := data.ComposeTransferDetail(tx, snowflakeId)
			if err != nil {
				data.Rollback(tx)
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取提现详情失败: %v", err)})
			}

			transferDetailList = append(transferDetailList, transferDetail)
			batchesRequest, err := data.ComposeBatches(transferDetailList)
			if err != nil {
				return fmt.Errorf("组合批次失败: %v", err)
			}

			batchesResponse, err := pay.Batches(batchesRequest)
			if err != nil {
				return fmt.Errorf("发起批次失败: %v", err)
			}

			err = data.CreatePay(tx, batchesRequest.TotalAmount, batchesRequest.TotalNum, batchesResponse)
			if err != nil {
				data.Rollback(tx)
				return fmt.Errorf("创建支付记录失败: %v", err)
			}

			err = data.UpdateWithdrawalBatchId(tx, transferDetailList, batchesResponse)
			if err != nil {
				data.Rollback(tx)
				return fmt.Errorf("更新提现记录失败: %v", err)
			}
		}
	}

	data.Commit(tx)

	return c.JSON(tools.Response{
		Code:    0,
		Message: "审批提现成功",
		Result:  struct{}{},
	})
}
