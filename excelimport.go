package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/xuri/excelize/v2"
)

func ExcelImport(c fiber.Ctx) error {

	multipart, err := c.MultipartForm()
	if err != nil {
		return err
	}

	for _, file := range multipart.File["file"] {

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("upload/" + file.Filename)
		if err != nil {
			return err
		}

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// Close the file
		defer dst.Close()

		f, err := excelize.OpenFile("upload/" + file.Filename)
		if err != nil {
			return err
		}
		// 获取 Sheet1 上所有单元格
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, row := range rows[1:] {
			for colIndex, colCell := range row {

				if colIndex == 4 || colIndex == 5 {
					// 判断是否为数字
					if _, err := strconv.ParseInt(colCell, 10, 64); err != nil {
						return c.JSON(Resp{
							Code:    1,
							Message: "excel格式错误,请检查出货量或积分列",
							Result:  struct{}{},
						})
					}

					// 查询手机号是否存在
					snowflakeId, err := FindPhoneNumberContext(c.Context(), row[3])

					if err != nil {
						return c.JSON(Resp{
							Code:    1,
							Message: "查询手机号失败",
							Result:  struct{}{},
						})
					}

					if snowflakeId != 0 {

					} else {
						baseSQL := "INSERT INTO `users` (nike, phone, province, city, shipments, integral) VALUES ($1, $2, $3, $4, $5, $6)"

						db.ExecContext(c.Context(), baseSQL, row[0], row[1], row[2], row[3], row[4], row[5])
					}

				}
			}
		}

		// Remove the temporary file
		defer os.Remove(file.Filename)

	}
	// Send a string response to the client
	return c.SendString("Hello, World 👋!")
}
