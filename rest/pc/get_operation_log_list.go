package pc

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func GetOperationLogList(c fiber.Ctx) error {
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

	payload := new(entity.GetOperationLogListRequest)

	if payload.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if payload.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	payload.UserId = c.Query("user_id", "")

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	operationLog, err := data.GetOperationLogList(tx, payload)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取收支列表失败: %v", err)})
	}

	total, err := data.GetOperationLogCount(tx, payload)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取收支列表数量失败: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取收支列表成功",
		Result: map[string]any{
			"operation_log_list": operationLog,
			"total":              total,
		},
	})

}
