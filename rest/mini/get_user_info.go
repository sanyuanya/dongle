package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/tools"
)

func GetUserInfo(c fiber.Ctx) error {

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

	userDetail, err := data.GetUserDetailBySnowflakeID(snowflakeId)
	if err != nil {
		panic(fmt.Errorf("获取用户信息失败: %v", err))
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result:  userDetail,
	})
}
