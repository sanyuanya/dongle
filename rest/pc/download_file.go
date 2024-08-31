package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/tools"
)

func DownloadFile(c fiber.Ctx) error {
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

	filepath := c.Params("fileName", "")

	if filepath == "" {
		panic(tools.CustomError{Code: 40000, Message: "参数错误"})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("开启事务失败: %v", err)})
	}

	fileName, err := data.GetIncomeByPath(tx, "/api/pc/upload/"+filepath)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50004, Message: fmt.Sprintf("查询失败: %v", err)})
	}

	tx.Commit()

	// 设置响应头，以确保浏览器正确识别文件类型
	c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set(fiber.HeaderContentDisposition, "attachment; filename="+fileName)

	return c.SendFile("upload/" + filepath)

}
