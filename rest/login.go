package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

// mysecretpassword
func PcLogin(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()
	loginRequest := new(entity.LoginRequest)
	err := c.Bind().Body(loginRequest)
	if err != nil {
		panic(fmt.Errorf("请求参数错误 : %v", err))
	}

	snowflakeId, err := data.Login(loginRequest)

	if err != nil {
		panic(fmt.Errorf("登录失败 : %v", err))
	}

	token, err := tools.GenerateToken(snowflakeId, "admin")

	if err != nil {
		panic(fmt.Errorf("生成token失败 : %v", err))
	}

	err = data.SetApiToken(snowflakeId, token)

	if err != nil {
		panic(fmt.Errorf("设置token失败 : %v", err))
	}

	c.Response().Header.Set("Authorization", token)
	return c.JSON(Resp{
		Code:    0,
		Message: "登录成功",
		Result:  struct{}{},
	})
}
