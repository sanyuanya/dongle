package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

type SetUserInfoRequest struct {
	Nike        string `json:"nike"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	IDCard      string `json:"id_card"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	CompanyName string `json:"company_name"`
	Job         string `json:"job"`
	SnowflakeId int64  `json:"snowflake_id"`
}

func SetUserInfo(c fiber.Ctx) error {
	authorization := c.Get("Authorization")

	token, err := ValidateToken(authorization)

	if err != nil {
		return c.JSON(Resp{
			Code:    40001,
			Message: "Unauthorized",
			Result:  struct{}{},
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userID"].(int64)

	payload := new(SetUserInfoRequest)

	err = c.Bind().Body(payload)

	if err != nil {
		return c.JSON(Resp{
			Code:    10000,
			Message: "参数错误",
			Result:  struct{}{},
		})
	}

	payload.SnowflakeId = userID
	err = UpdateUserInfo(payload)
	if err != nil {
		return c.JSON(Resp{
			Code:    50000,
			Message: err.Error(),
			Result:  struct{}{},
		})
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "success",
		Result:  struct{}{},
	})
}

func GetUserInfo(c fiber.Ctx) error {

	authorization := c.Get("Authorization")

	token, err := ValidateToken(authorization)

	if err != nil {
		return c.JSON(Resp{
			Code:    40001,
			Message: "Unauthorized",
			Result:  struct{}{},
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userID"].(int64)

	userDetail, err := GetUserDetailBySnowflakeID(userID)
	if err != nil {
		return c.JSON(Resp{
			Code:    50000,
			Message: err.Error(),
			Result:  struct{}{},
		})
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "success",
		Result:  userDetail,
	})
}
