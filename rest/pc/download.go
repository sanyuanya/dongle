package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/tools"
	"github.com/xuri/excelize/v2"
)

func Download(c fiber.Ctx) error {
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

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	// 查询所有的产品信息
	product, err := data.GetProductAll(tx)
	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取产品列表失败: %v", err)})
	}
	// 生成excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	fileHeader := []string{
		"姓名",
		"省份",
		"地市",
		"手机号",
	}

	for _, v := range product {
		fileHeader = append(fileHeader, fmt.Sprintf("%s 出货量", v.Name))
	}

	for i, v := range fileHeader {
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'A'+i, 1), v)
	}

	f.SetColWidth("Sheet1", "A", "G", 20)
	f.SetRowHeight("Sheet1", 1, 30)

	f.SetActiveSheet(1)

	if err := f.SaveAs("客户导入模版.xlsx"); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("保存文件失败")
	}
	// 设置响应头，以确保浏览器正确识别文件类型
	c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set(fiber.HeaderContentDisposition, "attachment; filename=客户导入模版.xlsx")

	return c.SendFile("客户导入模版.xlsx")
}