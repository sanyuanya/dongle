package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func UpdateUserInfo(c fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(tools.Response{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	_, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	payload := new(entity.UpdateUserDetailRequest)
	err = c.Bind().Body(payload)
	if err != nil {
		panic(fmt.Errorf("无法绑定请求体: %v", err))
	}

	err = data.UpdateUserDetail(payload)

	if err != nil {
		panic(fmt.Errorf("无法更新用户信息: %v", err))
	}
	return c.JSON(tools.Response{
		Code:    0,
		Message: "更新用户信息成功",
		Result:  struct{}{},
	})
}
