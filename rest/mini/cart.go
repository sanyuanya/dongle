package mini

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func CartIndex(c fiber.Ctx) error {
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

	payload := new(entity.GetCartListRequest)

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}
	payload.SnowflakeId = snowflakeId

	if payload.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if payload.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	cartList, err := data.GetCartList(tx, payload)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取购物车列表失败: %v", err)})
	}

	total, err := data.CartListTotal(tx, payload)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取购物车列表总数失败: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取购物车列表成功",
		Result: map[string]any{
			"cartList": cartList,
			"total":    total,
		},
	})
}

func CartAdd(c fiber.Ctx) error {
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
	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}
	payload := new(entity.AddCardRequest)
	payload.UserId = snowflakeId
	if err := c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("参数错误: %v", err)})
	}
	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}
	payload.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	if err := data.AddCart(tx, payload); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加购物车失败: %v", err)})
	}
	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加购物车成功",
		Result:  struct{}{},
	})
}

func CartUpdate(c fiber.Ctx) error {
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
	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}
	payload := new(entity.UpdateCardRequest)
	payload.SnowflakeId = c.Params("cartId", "")
	payload.UserId = snowflakeId
	if err := c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("参数错误: %v", err)})
	}
	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}
	if id, err := data.FindByCartSnowflakeId(tx, payload.SnowflakeId); err != nil || id == "" {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取购物车信息失败: %v", err)})
	}
	if err := data.UpdateCart(tx, payload); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新购物车失败: %v", err)})
	}
	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "更新购物车成功",
		Result:  struct{}{},
	})
}

func CartDelete(c fiber.Ctx) error {
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
	userId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := new(entity.DeleteCardRequest)
	if err := c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("参数错误: %v", err)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	for _, itemId := range payload.CartIdList {
		if id, err := data.FindByCartSnowflakeId(tx, itemId); err != nil || id == "" {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取购物车信息失败: %v", err)})
		}
		if err := data.DeleteCart(tx, itemId, userId); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("删除购物车失败: %v", err)})
		}
	}

	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "删除购物车成功",
		Result:  struct{}{},
	})
}
