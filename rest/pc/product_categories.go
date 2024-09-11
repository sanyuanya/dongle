package pc

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func GetProductCategoriesList(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := &entity.GetProductCategoriesListRequest{}

	if payload.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if payload.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	payload.Keyword = c.Query("keyword", "")

	if payload.Status, err = strconv.ParseInt(c.Query("status", "0"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	productCategoriesList, err := data.GetProductCategoriesList(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取产品分类列表失败: %v", err)})
	}

	total, err := data.GetProductCategoriesListCount(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取产品分类总数失败: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取产品分类列表成功",
		Result: map[string]interface{}{
			"data":  productCategoriesList,
			"total": total,
		},
	})

}

func AddProductCategories(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	addProductCategoriesRequest := new(entity.AddProductCategoriesRequest)

	err = c.Bind().Body(addProductCategoriesRequest)

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	addProductCategoriesRequest.SnowflakeId = tools.SnowflakeUseCase.NextVal()

	if addProductCategoriesRequest.Status > 2 {
		panic(tools.CustomError{Code: 40000, Message: "status 参数错误"})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	existsId, err := data.FindByProductCategoriesName(tx, addProductCategoriesRequest.Name)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50004, Message: fmt.Sprintf("查询产品分类失败: %v", err)})
	}

	if existsId != "" {
		tx.Rollback()
		panic(tools.CustomError{Code: 50005, Message: "产品分类已存在"})
	}

	err = data.AddProductCategories(tx, addProductCategoriesRequest)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("添加产品分类失败: %v", err)})
	}

	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加产品分类成功",
		Result:  struct{}{},
	})

}

func UpdateProductCategories(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	updateProductCategoriesRequest := new(entity.UpdateProductCategoriesRequest)

	updateProductCategoriesRequest.SnowflakeId = c.Params("productCategoriesId", "")

	err = c.Bind().Body(updateProductCategoriesRequest)

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	if updateProductCategoriesRequest.Status > 2 {
		panic(tools.CustomError{Code: 40000, Message: "status 参数错误"})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	existsId, err := data.FindByProductCategoriesName(tx, updateProductCategoriesRequest.Name)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50004, Message: fmt.Sprintf("查询产品分类失败: %v", err)})
	}

	if existsId != "" {
		tx.Rollback()
		panic(tools.CustomError{Code: 50005, Message: "产品分类已存在"})
	}

	err = data.UpdateProductCategories(tx, updateProductCategoriesRequest)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50009, Message: fmt.Sprintf("更新产品分类失败: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "更新产品分类成功",
		Result:  struct{}{},
	})

}

func DeleteProductCategories(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	productCategoriesId := c.Params("productCategoriesId", "")

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	err = data.DeleteProductCategories(tx, productCategoriesId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50008, Message: fmt.Sprintf("删除产品分类失败: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "删除产品分类成功",
		Result:  struct{}{},
	})

}
