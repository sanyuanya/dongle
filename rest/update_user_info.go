package rest

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
	"github.com/sanyuanya/dongle/wechat"
)

func UpdateUserInfo(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	_ = snowflakeId
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	payload := new(entity.UpdateUserInfoRequest)

	err = c.Bind().Body(payload)

	if err != nil {
		panic(fmt.Errorf("参数错误: %v", err))
	}

	accessToken, expiresIn, err := data.FindAccessTokenByAppId()

	if err != nil {
		panic(fmt.Errorf("获取appId失败: %v", err))
	}

	if expiresIn-30 <= time.Now().Unix() {
		getAccessTokenResp, err := wechat.GetAccessToken()
		if err != nil {
			panic(fmt.Errorf("获取access_token失败: %v", err))
		}

		accessToken = getAccessTokenResp.AccessToken
		expiresIn = time.Now().Unix() + getAccessTokenResp.ExpiresIn
		data.UpdateAccessTokenAndExpiresIn(accessToken, expiresIn)
	}

	getPhoneNumberResp, err := wechat.GetPhoneNumber(payload.Code, accessToken)

	if err != nil {
		panic(fmt.Errorf("获取用户手机号失败: %v", err))
	}

	if getPhoneNumberResp.Errcode != 0 {
		panic(fmt.Errorf("获取用户手机号失败: %v", getPhoneNumberResp.Errmsg))
	}

	userInfo := &entity.UserInfo{
		Nick:        payload.Nick,
		Avatar:      payload.Avatar,
		Phone:       getPhoneNumberResp.PhoneInfo.PhoneNumber,
		SnowflakeId: snowflakeId,
	}

	err = data.UpdateUser(userInfo)
	if err != nil {
		panic(fmt.Errorf("修改用户信息失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "成功",
		Result:  struct{}{},
	})
}
