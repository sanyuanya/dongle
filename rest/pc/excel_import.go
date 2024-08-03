package pc

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
	"github.com/xuri/excelize/v2"
)

func ExcelImport(c fiber.Ctx) error {

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

	multipart, err := c.MultipartForm()

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	batch := tools.SnowflakeUseCase.NextVal()

	for _, file := range multipart.File["file"] {

		// Remove the temporary file
		defer os.Remove("upload/" + file.Filename)

		src, err := file.Open()
		if err != nil {
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法打开文件: %v", err)})
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("upload/" + file.Filename)
		if err != nil {
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法创建文件: %v", err)})
		}

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法复制文件: %v", err)})
		}

		// Close the file
		defer dst.Close()

		f, err := excelize.OpenFile("upload/" + file.Filename)
		if err != nil {
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法打开文件: %v", err)})
		}
		// 获取 Sheet1 上所有单元格
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法获取行: %v", err)})
		}

		for rowIndex, row := range rows[1:] {

			importUserInfo := new(entity.ImportUserInfo)

			for colIndex, colCell := range row {

				if colIndex == 4 || colIndex == 5 {
					// 判断是否为数字
					if _, err := strconv.ParseInt(colCell, 10, 64); err != nil {
						panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 第 %d 列, 格式错误: %v", rowIndex+1, colIndex+1, err)})
					}
				}

				importUserInfo.Nick = row[0]
				importUserInfo.Province = row[1]
				importUserInfo.City = row[2]
				importUserInfo.Phone = row[3]

				// 更新用户积分和出货量
				importUserInfo.Shipments, err = strconv.ParseInt(row[4], 10, 64)
				if err != nil {
					panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 第 %d 列, 格式错误: %v", rowIndex+1, colIndex+1, err)})
				}

				importUserInfo.Integral, err = strconv.ParseInt(row[5], 10, 64)
				if err != nil {
					panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 第 %d 列, 格式错误: %v", rowIndex+1, colIndex+1, err)})
				}
			}

			// 查询手机号是否存在
			snowflakeId, err := data.FindPhoneNumberContext(row[3])

			if err != nil {
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("查询手机号失败: %v", err)})
			}

			if snowflakeId != 0 {
				err := data.UpdateUserIntegralAndShipments(snowflakeId, importUserInfo.Integral, importUserInfo.Shipments)
				if err != nil {
					panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("更新用户积分和出货量失败: %v", err)})
				}
			} else {
				// 新增用户
				importUserInfo.SnowflakeId = tools.SnowflakeUseCase.NextVal()
				err := data.ImportUserInfo(importUserInfo)

				if err != nil {
					panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("新增用户失败: %v", err)})
				}
			}

			addIncomeExpenseRequest := new(entity.AddIncomeExpenseRequest)

			addIncomeExpenseRequest.SnowflakeId = tools.SnowflakeUseCase.NextVal()
			addIncomeExpenseRequest.Summary = "分红奖励"
			addIncomeExpenseRequest.Integral = importUserInfo.Integral
			addIncomeExpenseRequest.Shipments = importUserInfo.Shipments
			addIncomeExpenseRequest.UserId = snowflakeId
			addIncomeExpenseRequest.Batch = batch

			err = data.AddIncomeExpense(addIncomeExpenseRequest)

			if err != nil {
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("新增收支记录失败: %v", err)})
			}

		}

	}
	// Send a string response to the client
	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})
}
