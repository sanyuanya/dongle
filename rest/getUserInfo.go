package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
)

func SetUserInfo(c fiber.Ctx) error {
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

	if role != "user" {
		panic(fmt.Errorf("unauthorized"))
	}

	payload := new(entity.SetUserInfoRequest)

	payload.SnowflakeId, err = strconv.ParseInt(snowflakeIdStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err))
	}

	err = c.Bind().Body(payload)

	if err != nil {
		panic(fmt.Errorf("参数错误: %v", err))
	}

	err = data.UpdateUserInfo(payload)
	if err != nil {
		panic(fmt.Errorf("更新用户信息失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})
}

func GetUserInfo(c fiber.Ctx) error {

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

		snowflakeId, err := strconv.ParseInt(c.Query("snowflake_id"), 10, 64)
		if err != nil {
			panic(fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err))
		}

		userDetail, err := data.GetUserDetailBySnowflakeID(snowflakeId)

		_ = userDetail
		if err != nil {
			panic(fmt.Errorf("获取用户信息失败: %v", err))
		}

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

	if role != "user" {
		panic(fmt.Errorf("unauthorized"))
	}

	snowflakeId, err := strconv.ParseInt(snowflakeIdStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err))
	}

	userDetail, err := data.GetUserDetailBySnowflakeID(snowflakeId)
	if err != nil {
		panic(fmt.Errorf("获取用户信息失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "success",
		Result:  userDetail,
	})
}
