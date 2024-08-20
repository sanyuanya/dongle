package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/tools"
)

func GetProductGroupList(c fiber.Ctx) error {
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

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	productGroupList, err := data.GetProductGroupList(tx, snowflakeId)
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取产品组列表失败: %v", err)})
	}

	var totalIntegral int64

	var totalShipment int64
	// 计算总积分
	for _, item := range productGroupList {
		totalIntegral += item.Merge
		totalShipment += item.Shipments
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取产品组列表成功",
		Result: map[string]any{
			"productGroupList": productGroupList,
			"total":            len(productGroupList),
			"total_integral":   totalIntegral,
			"total_shipments":  totalShipment,
		},
	})

}
