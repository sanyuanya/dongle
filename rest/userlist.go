package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
)

func UserList(c fiber.Ctx) error {

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

	userPageListRequest := &entity.UserPageListRequest{}

	if userPageListRequest.Page, err = strconv.ParseInt(c.Query("page"), 10, 64); err != nil {
		panic(fmt.Errorf("page 参数错误: %v", err))
	}

	if userPageListRequest.PageSize, err = strconv.ParseInt(c.Query("page_size"), 10, 64); err != nil {
		panic(fmt.Errorf("page_size 参数错误: %v", err))
	}

	if userPageListRequest.IsWhite, err = strconv.ParseInt(c.Query("is_white"), 10, 64); err != nil {
		panic(fmt.Errorf("is_white 参数错误: %v", err))
	}

	userPageListRequest.Keyword = c.Query("keyword")

	userList, err := data.GetUserPageList(userPageListRequest)
	if err != nil {
		panic(fmt.Errorf("获取用户列表失败: %v", err))
	}

	userTotal, err := data.GetUserPageCount(userPageListRequest)
	if err != nil {
		panic(fmt.Errorf("获取用户总数失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "success",
		Result: map[string]any{
			"data":  userList,
			"total": userTotal,
		},
	})
}
