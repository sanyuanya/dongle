package pc

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/entity"
	expressdelivery "github.com/sanyuanya/dongle/express_delivery"
	"github.com/sanyuanya/dongle/tools"
)

func PollQuery(c fiber.Ctx) error {
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
	payload := &entity.PollQueryRequest{}
	if err := c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	rdb := tools.Redis{}
	if err := rdb.NewClient(); err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Redis 客户端: %v", err)})
	}

	value, err := rdb.GetLogisticsInformation(payload.Num)
	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法获取 Redis 客户端数据: %v", err)})
	}
	if value == "" {
		resp, err := expressdelivery.PollQuery(payload)
		if err != nil {
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("失败: %v", err)})
		}
		information, err := json.Marshal(resp)
		if err != nil {
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("序列化物流信息失败: %v", err)})
		}
		if err := rdb.SetLogisticsInformation(payload.Num, string(information)); err != nil {
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("设置订单物流信息失败: %v", err)})
		}
		return c.JSON(resp)
	} else {
		var resp *entity.PollQueryResponse
		if err := json.Unmarshal([]byte(value), &resp); err != nil {
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("反序列化物流信息失败: %v", err)})
		}
		return c.JSON(resp)
	}
}
