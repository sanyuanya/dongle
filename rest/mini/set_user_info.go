package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func SetUserInfo(c fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(tools.Response{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")

	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	payload := new(entity.SetUserInfoRequest)

	payload.SnowflakeId = snowflakeId

	err = c.Bind().Body(payload)

	if err != nil {
		panic(fmt.Errorf("参数错误: %v", err))
	}

	err = data.UpdateUserInfo(payload)
	if err != nil {
		panic(fmt.Errorf("更新用户信息失败: %v", err))
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})
}
