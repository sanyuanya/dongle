package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func AddWhite(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(tools.Response{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	addWhiteList := new(entity.AddWhiteRequest)
	err = c.Bind().Body(addWhiteList)

	if err != nil {
		panic(fmt.Errorf("参数错误: %v", err))
	}

	err = data.AddWhite(addWhiteList)

	if err != nil {
		panic(fmt.Errorf("添加白名单失败: %v", err))
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加白名单成功",
		Result:  struct{}{},
	})

}
