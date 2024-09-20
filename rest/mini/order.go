package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func Submit(c fiber.Ctx) error {
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
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := new(entity.SubmitOrderRequest)

	if err = c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("参数错误: %v", err)})
	}

	rdb := tools.Redis{}

	if err = rdb.NewClient(); err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Redis 客户端: %v", err)})
	}

	defer rdb.Client.Close()

	if result, err := rdb.DeductStock(payload.SkuId, payload.Quantity); err != nil || !result {
		panic(tools.CustomError{Code: 50002, Message: fmt.Sprintf("库存不足,下单失败")})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	if err = data.UpdateSkuStockQuantity(tx, payload.CommodityId, payload.SkuId, payload.Quantity); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新库存失败: %v", err)})
	}

	if err = data.AddOrder(tx, payload); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加订单失败: %v", err)})
	}

	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "下单成功",
		Result:  struct{}{},
	})
}
