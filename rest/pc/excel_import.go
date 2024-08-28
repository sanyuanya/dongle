package pc

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

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

	values := multipart.Value

	startTime := values["start_time"]
	endTime := values["end_time"]

	if len(startTime) == 0 || len(endTime) == 0 {
		panic(tools.CustomError{Code: 40000, Message: "开始时间和结束时间不能为空"})
	}

	beginTime, err := tools.ValidateTimestamp(startTime[0])

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("开始时间格式错误: %v", err)})
	}

	finishTime, err := tools.ValidateTimestamp(endTime[0])

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("结束时间格式错误: %v", err)})
	}

	if beginTime.After(finishTime) {
		panic(tools.CustomError{Code: 40000, Message: "开始时间不能晚于结束时间"})
	}

	if finishTime.After(time.Now()) {
		panic(tools.CustomError{Code: 40000, Message: "结束时间不能晚于当前时间"})
	}

	// 查询当前日期是否已经导入
	exist, err := data.CheckImportedAt(beginTime, finishTime)

	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("查询导入日期失败: %v", err)})
	}

	if exist {
		panic(tools.CustomError{Code: 40000, Message: "当前日期已经导入, 请勿重复导入"})
	}

	file := multipart.File["file"][0]

	// Remove the temporary file
	defer os.Remove("upload/" + file.Filename)

	src, err := file.Open()
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法打开文件: %v", err)})
	}
	defer src.Close()

	err = os.MkdirAll("upload", os.ModePerm)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法创建文件夹: %v", err)})
	}

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

	if rows[0][0] != "日期" || rows[0][1] != "姓名" || rows[0][2] != "省份" || rows[0][3] != "地市" || rows[0][4] != "手机号" {
		panic(tools.CustomError{Code: 40000, Message: "表头错误"})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("开启事务失败: %v", err)})
	}

	layout := "2006-01-02"
	for rowIndex, row := range rows[1:] {
		length := len(row)

		if length <= 5 {
			tx.Rollback()
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 列数错误", rowIndex+1)})
		}

		importUserInfo := new(entity.ImportUserInfo)

		importUserInfo.ImportdAt, err = time.Parse(layout, row[0])

		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 日期格式错误, 格式为: 年-月-日 例如 2024-08-07", rowIndex+1)})
		}

		importUserInfo.Nick = row[1]
		importUserInfo.Province = row[2]
		importUserInfo.City = row[3]
		importUserInfo.Phone = row[4]

		importUserInfo.WithdrawablePoints, err = strconv.ParseInt(row[length-1], 10, 64)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 可提现积分格式错误", rowIndex+1)})
		}

		if len(importUserInfo.Phone) != 11 {
			tx.Rollback()
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 手机号错误", rowIndex+1)})
		}

		for colIndex, colCell := range row[5 : length-1] {

			shipment, err := strconv.ParseInt(colCell, 10, 64)
			if err != nil {
				tx.Rollback()
				panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 第 %d 列, 单元格: %v 格式错误", rowIndex+2, colIndex+5, colCell)})
			}

			if shipment < 0 || shipment > 100000 {
				tx.Rollback()
				panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("第 %d 行, 第 %d 列, 单元格: %v 不能为负数、或大于 10 万", rowIndex+2, colIndex+5, colCell)})
			}

			if shipment == 0 {
				continue
			}

			productName := strings.TrimSpace(strings.ReplaceAll(rows[0][colIndex+5], "出货量", ""))
			product, err := data.FindProductByName(tx, productName)
			if err != nil {
				tx.Rollback()
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("查询产品失败: %v", err)})
			}

			if product == nil {
				tx.Rollback()
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("第 %d 行, 第 %d 列, 产品不存在", rowIndex+1, colIndex+5)})
			}

			// 统计用户积分
			importUserInfo.Integral = shipment * product.Integral

			// 查询手机号是否存在
			snowflakeId, err := data.FindPhoneNumberContext(tx, row[4])

			if err != nil {
				tx.Rollback()
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("查询手机号失败: %v", err)})
			}

			if snowflakeId != "" {
				err := data.UpdateUserIntegralAndShipments(tx, snowflakeId, importUserInfo.Integral, importUserInfo.Shipments, importUserInfo.WithdrawablePoints)
				if err != nil {
					tx.Rollback()
					panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("更新用户积分和出货量失败: %v", err)})
				}
			} else {
				// 新增用户
				importUserInfo.SnowflakeId = tools.SnowflakeUseCase.NextVal()
				err := data.ImportUserInfo(tx, importUserInfo)

				if err != nil {
					tx.Rollback()
					panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("新增用户失败: %v", err)})
				}

				snowflakeId = importUserInfo.SnowflakeId
			}

			addIncomeExpenseRequest := new(entity.AddIncomeExpenseRequest)
			addIncomeExpenseRequest.SnowflakeId = tools.SnowflakeUseCase.NextVal()
			addIncomeExpenseRequest.Summary = "分红奖励"
			addIncomeExpenseRequest.Integral = importUserInfo.Integral
			addIncomeExpenseRequest.Shipments = shipment
			addIncomeExpenseRequest.UserId = snowflakeId
			addIncomeExpenseRequest.Batch = batch
			addIncomeExpenseRequest.ProductId = product.SnowflakeId
			addIncomeExpenseRequest.ProductIntegral = product.Integral
			addIncomeExpenseRequest.ImportdAt = importUserInfo.ImportdAt

			err = data.AddIncomeExpense(tx, addIncomeExpenseRequest)
			if err != nil {
				tx.Rollback()
				panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("新增收支记录失败: %v", err)})
			}
		}
	}

	tx.Commit()
	// Send a string response to the client
	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})
}
