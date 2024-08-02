package rest

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/snowflake"
	"github.com/sanyuanya/dongle/tools"
	"github.com/sanyuanya/dongle/wechat"
)

func MiniLogin(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	miniLoginRequest := new(entity.MiniLoginRequest)

	err := c.Bind().Body(miniLoginRequest)
	if err != nil {
		panic(fmt.Errorf("请求参数错误 : %v", err))
	}

	code2SessionResp, err := wechat.Code2Session(miniLoginRequest.JsCode)
	if err != nil {
		panic(fmt.Errorf("获取openid失败: %v", err))
	}

	if code2SessionResp.ErrCode != 0 {
		panic(fmt.Errorf("获取openid失败: %v", code2SessionResp.ErrMsg))
	}

	snowflakeId, err := data.FindOpenId(code2SessionResp.OpenID)

	if err != nil {
		panic(fmt.Errorf("openid查询失败: %v", err))
	}

	if snowflakeId == 0 {

		snowflakeId = snowflake.SnowflakeUseCase.NextVal()

		registerUserRequest := &entity.RegisterUserRequest{
			SnowflakeId: snowflakeId,
			OpenId:      code2SessionResp.OpenID,
			SessionKey:  code2SessionResp.SessionKey,
		}

		err := data.RegisterUser(registerUserRequest)
		if err != nil {
			panic(fmt.Errorf("openid注册失败: %v", err))
		}
	} else {
		err := data.UpdateSessionKey(code2SessionResp.OpenID, code2SessionResp.SessionKey)
		if err != nil {
			panic(fmt.Errorf("openid更新失败: %v", err))
		}
	}

	token, err := tools.GenerateToken(snowflakeId, "user")

	if err != nil {
		panic(fmt.Errorf("生成token失败: %v", err))
	}

	c.Response().Header.Set("Authorization", token)

	return c.JSON(Resp{
		Code:    0,
		Message: "成功",
		Result: map[string]any{
			"openid":       code2SessionResp.OpenID,
			"session_key":  code2SessionResp.SessionKey,
			"snowflake_id": snowflakeId,
		},
	})
}

func Register(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	req := new(entity.RegisterRequest)

	err := c.Bind().Body(req)

	if err != nil {
		panic(fmt.Errorf("参数错误: %v", err))
	}

	code2SessionResp, err := wechat.Code2Session(req.JsCode)

	if err != nil {
		panic(fmt.Errorf("获取openid失败: %v", err))
	}

	if code2SessionResp.ErrCode != 0 {
		panic(fmt.Errorf("获取openid失败: %v", code2SessionResp.ErrMsg))
	}

	accessToken, expiresIn, err := data.FindAppId()

	if err != nil {
		panic(fmt.Errorf("获取appId失败: %v", err))
	}

	if expiresIn-30 <= time.Now().Unix() {
		// getAccessToken
		getAccessTokenResp, err := wechat.GetAccessToken()
		if err != nil {
			panic(fmt.Errorf("获取access_token失败: %v", err))
		}

		accessToken = getAccessTokenResp.AccessToken
		expiresIn = time.Now().Unix() + getAccessTokenResp.ExpiresIn
		data.UpdateAccessTokenAndExpiresIn(accessToken, expiresIn)
	}

	getPhoneNumberResp, err := wechat.GetPhoneNumber(req.Code, accessToken)

	if err != nil {
		panic(fmt.Errorf("获取用户手机号失败: %v", err))
	}

	if getPhoneNumberResp.Errcode != 0 {
		panic(fmt.Errorf("获取用户手机号失败: %v", getPhoneNumberResp.Errmsg))
	}

	snowflakeId, err := data.FindPhoneNumberContext(getPhoneNumberResp.PhoneInfo.PhoneNumber)

	if err != nil {
		panic(fmt.Errorf("手机号查询失败: %v", err))
	}

	userInfo := &entity.UserInfo{
		OpenID:     code2SessionResp.OpenID,
		Nick:       req.Nick,
		Avatar:     req.Avatar,
		Phone:      getPhoneNumberResp.PhoneInfo.PhoneNumber,
		SessionKey: code2SessionResp.SessionKey,
	}

	if snowflakeId != 0 {
		userInfo.SnowflakeId = snowflakeId
		err := data.UpdateUser(userInfo)
		if err != nil {
			panic(fmt.Errorf("用户修改失败: %v", err))
		}
	} else {
		userInfo.SnowflakeId = snowflake.SnowflakeUseCase.NextVal()
		err := data.RegisterUser(userInfo)
		if err != nil {
			panic(fmt.Errorf("用户注册失败: %v", err))
		}
	}

	apiToken, err := GenerateToken(userInfo.SnowflakeId, "user")

	if err != nil {
		panic(fmt.Errorf("生成token失败: %v", err))
	}

	c.Response().Header.Set("Authorization", apiToken)

	return c.JSON(Resp{
		Code:    0,
		Message: "成功",
		Result: map[string]any{
			"openid":       code2SessionResp.OpenID,
			"session_key":  code2SessionResp.SessionKey,
			"phone":        getPhoneNumberResp.PhoneInfo.PhoneNumber,
			"snowflake_id": userInfo.SnowflakeId,
		},
	})
}
