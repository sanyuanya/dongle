package main

import "github.com/gofiber/fiber/v3"

type RegisterRequest struct {
	Code   string `json:"code"`
	JsCode string `json:"js_code"`
	Nike   string `json:"nike"`
	Avatar string `json:"avatar"`
}

type UserInfo struct {
	OpenID      string `json:"open_id"`
	Nike        string `json:"nike"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	ApiToken    string `json:"api_token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	SessionKey  string `json:"session_key"`
	SnowflakeId int64  `json:"snowflake_id"`
}

func Register(c fiber.Ctx) error {
	req := new(RegisterRequest)

	err := c.Bind().Body(req)

	if err != nil {
		return c.JSON(Resp{
			Code:    10000,
			Message: "参数错误",
			Result:  struct{}{},
		})
	}

	code2SessionResp, err := Code2Session(req.JsCode)

	if err != nil {
		return c.JSON(Resp{
			Code:    10001,
			Message: "获取openid失败",
			Result:  struct{}{},
		})
	}

	if code2SessionResp.ErrCode != 0 {
		return c.JSON(Resp{
			Code:    code2SessionResp.ErrCode,
			Message: code2SessionResp.ErrMsg,
			Result:  struct{}{},
		})
	}

	// getAccessToken
	getAccessTokenResp, err := GetAccessToken()
	if err != nil {
		return c.JSON(Resp{
			Code:    10002,
			Message: "获取access_token失败",
			Result:  struct{}{},
		})
	}

	getPhoneNumberResp, err := GetPhoneNumber(req.Code, getAccessTokenResp.AccessToken)

	if err != nil {
		return c.JSON(Resp{
			Code:    10003,
			Message: "获取用户手机号失败",
			Result:  struct{}{},
		})
	}

	if getPhoneNumberResp.Errcode != 0 {
		return c.JSON(Resp{
			Code:    getPhoneNumberResp.Errcode,
			Message: getPhoneNumberResp.Errmsg,
			Result:  struct{}{},
		})
	}

	snowflakeId, err := FindPhoneNumberContext(c.Context(), getPhoneNumberResp.PhoneInfo.PhoneNumber)

	if err != nil {
		return c.JSON(Resp{
			Code:    10004,
			Message: "获取用户手机号失败",
			Result:  struct{}{},
		})
	}

	userInfo := &UserInfo{
		OpenID:      code2SessionResp.OpenID,
		Nike:        req.Nike,
		Avatar:      req.Avatar,
		Phone:       getPhoneNumberResp.PhoneInfo.PhoneNumber,
		AccessToken: getAccessTokenResp.AccessToken,
		ExpiresIn:   getAccessTokenResp.ExpiresIn,
		SessionKey:  code2SessionResp.SessionKey,
	}

	if snowflakeId != 0 {
		userInfo.SnowflakeId = snowflakeId
		err := UpdateUser(userInfo)
		if err != nil {
			return c.JSON(Resp{
				Code:    10005,
				Message: "修改用户失败",
				Result:  struct{}{},
			})
		}
	} else {
		userInfo.SnowflakeId = snowflake.NextVal()
		err := RegisterUser(userInfo)

		if err != nil {
			return c.JSON(Resp{
				Code:    10005,
				Message: "注册用户失败",
				Result:  struct{}{},
			})
		}
	}

	apiToken, err := GenerateToken(userInfo.SnowflakeId)

	if err != nil {
		return c.JSON(Resp{
			Code:    10006,
			Message: "生成token失败",
			Result:  struct{}{},
		})
	}

	c.Response().Header.Set("api_token", apiToken)

	return c.JSON(Resp{
		Code:    0,
		Message: "成功",
		Result: map[string]any{
			"openid":       code2SessionResp.OpenID,
			"session_key":  code2SessionResp.SessionKey,
			"phoneNumber":  getPhoneNumberResp.PhoneInfo.PhoneNumber,
			"api_token":    apiToken,
			"snowflake_id": userInfo.SnowflakeId,
		},
	})

}
