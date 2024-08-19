package pc

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func WithdrawalList(c fiber.Ctx) error {

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

	_, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	withdrawalPageListRequest := &entity.WithdrawalPageListRequest{}

	if withdrawalPageListRequest.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if withdrawalPageListRequest.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	if withdrawalPageListRequest.LifeCycle, err = strconv.ParseInt(c.Query("life_cycle", "0"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("life_cycle 参数错误: %v", err)})
	}

	withdrawalPageListRequest.Date = c.Query("date")

	withdrawalPageListRequest.Keyword = c.Query("keyword")

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	withdrawalList, err := data.WithdrawalPageList(tx, withdrawalPageListRequest)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取提现列表失败: %v", err)})
	}

	total, err := data.WithdrawalListCount(tx, withdrawalPageListRequest)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取提现总数失败: %v", err)})
	}

	data.Commit(tx)
	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result: map[string]any{
			"data":  withdrawalList,
			"total": total,
		},
	})

}
