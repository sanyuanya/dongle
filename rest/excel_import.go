package rest

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/snowflake"
	"github.com/sanyuanya/dongle/tools"
	"github.com/xuri/excelize/v2"
)

func ExcelImport(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	multipart, err := c.MultipartForm()

	if err != nil {
		panic(err)
	}

	batch := snowflake.SnowflakeUseCase.NextVal()

	for _, file := range multipart.File["file"] {

		src, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("upload/" + file.Filename)
		if err != nil {
			panic(err)
		}

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			panic(err)
		}

		// Close the file
		defer dst.Close()

		f, err := excelize.OpenFile("upload/" + file.Filename)
		if err != nil {
			panic(err)
		}
		// 获取 Sheet1 上所有单元格
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			panic(err)
		}

		for rowIndex, row := range rows[1:] {

			importUserInfo := new(entity.ImportUserInfo)

			for colIndex, colCell := range row {

				if colIndex == 4 || colIndex == 5 {
					// 判断是否为数字
					if _, err := strconv.ParseInt(colCell, 10, 64); err != nil {

						panic(fmt.Errorf("第 %d 行, 第 %d 列, 格式错误: %v", rowIndex+1, colIndex+1, err))
					}
				}

				importUserInfo.Nick = row[0]
				importUserInfo.Province = row[1]
				importUserInfo.City = row[2]
				importUserInfo.Phone = row[3]

				// 更新用户积分和出货量
				importUserInfo.Shipments, err = strconv.ParseInt(row[4], 10, 64)
				if err != nil {
					panic(fmt.Errorf("出货量格式错误: %v", err))
				}

				importUserInfo.Integral, err = strconv.ParseInt(row[5], 10, 64)
				if err != nil {
					panic(fmt.Errorf("积分格式错误: %v", err))
				}
			}

			// 查询手机号是否存在
			snowflakeId, err := data.FindPhoneNumberContext(row[3])

			if err != nil {
				panic(fmt.Errorf("查询手机号失败: %v", err))
			}

			if snowflakeId != 0 {
				err := data.UpdateUserIntegralAndShipments(snowflakeId, importUserInfo.Integral, importUserInfo.Shipments)
				if err != nil {
					panic(fmt.Errorf("更新用户积分和出货量失败: %v", err))
				}
			} else {
				// 新增用户
				importUserInfo.SnowflakeId = snowflake.SnowflakeUseCase.NextVal()
				err := data.ImportUserInfo(importUserInfo)

				if err != nil {
					panic(fmt.Errorf("新增用户失败: %v", err))
				}
			}

			addIncomeExpenseRequest := new(entity.AddIncomeExpenseRequest)

			addIncomeExpenseRequest.SnowflakeId = snowflake.SnowflakeUseCase.NextVal()
			addIncomeExpenseRequest.Summary = "分红奖励"
			addIncomeExpenseRequest.Integral = importUserInfo.Integral
			addIncomeExpenseRequest.Shipments = importUserInfo.Shipments
			addIncomeExpenseRequest.UserId = snowflakeId
			addIncomeExpenseRequest.Batch = batch

			err = data.AddIncomeExpense(addIncomeExpenseRequest)

			if err != nil {
				panic(fmt.Errorf("新增收支记录失败: %v", err))
			}

		}
		// Remove the temporary file
		defer os.Remove(file.Filename)
	}
	// Send a string response to the client
	return c.JSON(Resp{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})
}
