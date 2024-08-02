package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
)

func GetWithdrawalList(c fiber.Ctx) error {

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
		panic(fmt.Errorf("未经授权: %v", err))
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
		panic(fmt.Errorf("未经授权"))
	}

	snowflakeId, err := strconv.ParseInt(snowflakeIdStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err))
	}

	getWithdrawalListRequest := new(entity.GetWithdrawalListRequest)

	if getWithdrawalListRequest.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(fmt.Errorf("page 参数错误: %v", err))
	}

	if getWithdrawalListRequest.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(fmt.Errorf("page_size 参数错误: %v", err))
	}

	getWithdrawalListRequest.Date = c.Query("date", "")

	withdrawalList, err := data.GetWithdrawalListBySnowflakeId(snowflakeId, getWithdrawalListRequest)

	if err != nil {
		panic(fmt.Errorf("获取提现列表失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    20000,
		Message: "获取提现列表成功",
		Result:  withdrawalList,
	})

}
