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

	// 查询所有的产品信息
	product, err := data.GetProductAll()
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取产品列表失败: %v", err)})
	}

	// 生成excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("生成Excel发生错误: %v", err)})
		}
	}()

	fileHeader := []string{
		"日期 (格式：年-月-日，正确示例：2021-01-01，错误示例：2021-1-1、2021/1/1。错误格式将导致导入失败)",
		"姓名",
		"省份",
		"地市",
		"手机号",
	}

	for _, v := range product {
		fileHeader = append(fileHeader, fmt.Sprintf("%s 出货量", v.Name))
	}

	maxCol := len(fileHeader)
	for i, v := range fileHeader {
		f.SetCellValue("Sheet1", fmt.Sprintf("%c%d", 'A'+i, 1), v)
	}

	f.SetColWidth("Sheet1", "A", fmt.Sprintf("%c", 'A'+maxCol), 40)
	f.SetRowHeight("Sheet1", 1, 50)

	f.SetActiveSheet(1)

	styleId, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("创建样式失败: %v", err)})
	}

	f.SetRowStyle("Sheet1", 1, 1, styleId)

	err = f.SetDocProps(&excelize.DocProperties{
		Title:          "客户导入模版",
		Creator:        "Dongle",
		Category:       "客户导入",
		ContentStatus:  "Reviewed",
		Description:    "客户导入模版",
		Identifier:     "xlsx",
		Keywords:       "客户导入",
		LastModifiedBy: "Dongle",
		Revision:       "1",
		Subject:        "客户导入",
		Version:        "1.0",
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("设置文件属性失败")
	}

	if err := f.SaveAs("客户导入模版.xlsx"); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("保存文件失败")
	}
	// 设置响应头，以确保浏览器正确识别文件类型
	c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set(fiber.HeaderContentDisposition, "attachment; filename=客户导入模版.xlsx")

	return c.SendFile("客户导入模版.xlsx")
}
