package pc

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
)

func OrderCallback(c fiber.Ctx) error {

	param := c.FormValue("param")
	taskId := c.FormValue("taskId")
	sign := c.FormValue("sign")

	log.Printf("接受寄件下单回调接口 param: %s, taskId: %s, sign: %s", param, taskId, sign)

	orderCallback := new(entity.OrderCallback)
	if err := json.Unmarshal([]byte(param), orderCallback); err != nil {
		log.Printf("接受寄件下单回调接口 Json Unmarshal 失败 param: %s, taskId: %s, sign: %s", param, taskId, sign)

		return c.JSON(map[string]any{
			"result":     false,
			"returnCode": "500",
			"message":    "Invalid parameters",
		})
	}

	tx, err := data.Transaction()
	if err != nil {
		log.Printf("接受寄件下单回调接口失败 开启事物失败: %v", err)
		return c.JSON(map[string]any{
			"result":     false,
			"returnCode": "500",
			"message":    "开启事务失败",
		})
	}

	if err := data.UpdatedShippingByThirdOrderId(tx, orderCallback.Data); err != nil {
		log.Printf("接受寄件下单回调接口更新数据库失败 orderCallback: %v", orderCallback.Data)
		return c.JSON(map[string]any{
			"result":     false,
			"returnCode": "500",
			"message":    "更新寄件信息失败",
		})
	}

	tx.Commit()

	return c.JSON(map[string]any{
		"result":     true,
		"returnCode": "200",
		"message":    "successfully updated order callback.",
	})
}
