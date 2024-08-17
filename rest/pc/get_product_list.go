package pc

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func GetProductList(c fiber.Ctx) error {

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

	payload := &entity.GetProductListRequest{}

	if payload.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if payload.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	payload.Keyword = c.Query("keyword", "")

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	productList, err := data.GetProductList(tx, payload)
	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取商品列表失败: %v", err)})
	}

	productTotal, err := data.GetProductTotal(tx, payload)
	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取商品总数失败: %v", err)})
	}

	data.Commit(tx)

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取商品列表成功",
		Result: map[string]interface{}{
			"product_list": productList,
			"total":        productTotal,
		},
	})

}