package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
)

func AddWhite(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	authorization := c.Get("Authorization")

	token, err := ValidateToken(authorization)

	if err != nil {
		panic(fmt.Errorf("unauthorized: %v", err))
	}

	claims := token.Claims.(jwt.MapClaims)

	snowflakeIdStr, ok := claims["snowflake_id"].(string)
	if !ok {
		panic(fmt.Errorf("snowflake_id is not a string"))
	}

	role, ok := claims["role"].(string)
	if !ok {
		panic(fmt.Errorf("role is not a string"))
	}

	if role != "admin" {
		panic(fmt.Errorf("unauthorized"))
	}

	snowflakeId, err := strconv.ParseInt(snowflakeIdStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err))
	}
	_ = snowflakeId

	addWhiteList := new(entity.AddWhiteRequest)
	err = c.Bind().Body(addWhiteList)

	if err != nil {
		panic(fmt.Errorf("参数错误: %v", err))
	}

	err = data.AddWhite(addWhiteList)

	if err != nil {
		panic(fmt.Errorf("添加白名单失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "添加白名单成功",
		Result:  struct{}{},
	})

}
