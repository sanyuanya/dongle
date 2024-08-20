package pc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
	"github.com/xuri/excelize/v2"
)

func ExportUser(c fiber.Ctx) error {

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

	// _, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	// if err != nil {
	// 	panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	// }

	exportUserRequest := &entity.ExportUserRequest{}

	var err error
	if exportUserRequest.IsWhite, err = strconv.ParseInt(c.Query("is_white", "0"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("is_white 参数错误: %v", err)})
	}

	exportUserRequest.Keyword = c.Query("keyword")

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	userList, err := data.GetUserList(tx, exportUserRequest)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取用户列表失败: %v", err)})
	}

	tx.Commit()

	// 生成excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("生成Excel发生错误: %v", err)})
		}
	}()

	f.SetCellValue("Sheet1", "A1", "姓名")
	f.SetCellValue("Sheet1", "B1", "省份")
	f.SetCellValue("Sheet1", "C1", "地市")
	f.SetCellValue("Sheet1", "D1", "手机号")
	f.SetCellValue("Sheet1", "E1", "积分")
	f.SetCellValue("Sheet1", "F1", "公司名称")
	f.SetCellValue("Sheet1", "G1", "职称")
	f.SetCellValue("Sheet1", "H1", "是否白名单")
	f.SetCellValue("Sheet1", "I1", "可提现积分")

	f.SetRowHeight("Sheet1", 1, 40)
	f.SetColWidth("Sheet1", "A", "C", 20)

	f.SetColWidth("Sheet1", "C", "I", 40)

	styleId, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("创建样式失败: %v", err)})
	}

	f.SetRowStyle("Sheet1", 1, 1, styleId)
	for i, user := range userList {

		var white string
		if user.IsWhite == 1 {
			white = "是"
		} else {
			white = "否"
		}
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'A', i+2), user.Nick)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'B', i+2), user.Province)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'C', i+2), user.City)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'D', i+2), user.Phone)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'E', i+2), user.Integral)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'F', i+2), user.CompanyName)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'G', i+2), user.Job)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'H', i+2), white)
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'I', i+2), user.WithdrawablePoints)
	}

	fileName := fmt.Sprintf("客户信息%s.xlsx", time.Now().Format("2006-01-02-150405"))

	if err := f.SaveAs(fileName); err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("保存文件失败: %v", err)})
	}
	// 设置响应头，以确保浏览器正确识别文件类型
	c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", fileName))

	return c.SendFile(fileName)
}
