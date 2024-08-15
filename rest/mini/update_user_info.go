package mini

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

			var code int
			var message string

			switch e := err.(type) {
			case tools.CustomError:
				code = e.Code
				message = e.Message
			case error:
				code = 50001
				message = e.Error()
			default:
				code = 50002
				message = fmt.Sprintf("%v", e)
			}

			c.JSON(tools.Response{
				Code:    code,
				Message: message,
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := new(entity.UpdateUserInfoRequest)

	err = c.Bind().Body(payload)

	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	accessToken, expiresIn, err := data.FindAccessTokenByAppId(tx)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取appId失败: %v", err)})
	}

	if expiresIn-30 <= time.Now().Unix() {
		getAccessTokenResp, err := wechat.GetAccessToken()
		if err != nil {
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取access_token失败: %v", err)})
		}

		accessToken = getAccessTokenResp.AccessToken
		expiresIn = time.Now().Unix() + getAccessTokenResp.ExpiresIn
		err = data.UpdateAccessTokenAndExpiresIn(tx, accessToken, expiresIn)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("更新access_token失败: %v", err)})
		}
	}

	getPhoneNumberResp, err := wechat.GetPhoneNumber(payload.Code, accessToken)

	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取用户手机号失败: %v", err)})
	}

	if getPhoneNumberResp.Errcode != 0 {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取用户手机号失败: %v", getPhoneNumberResp.Errmsg)})
	}

	userInfo := &entity.UserInfo{
		Nick:        payload.Nick,
		Avatar:      payload.Avatar,
		Phone:       getPhoneNumberResp.PhoneInfo.PhoneNumber,
		SnowflakeId: snowflakeId,
	}

	// 查询 手机号是否已经存在了
	userInfoReplace, err := data.FindUserByPhone(tx, userInfo.Phone)
	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("查询用户失败: %v", err)})
	}

	if userInfoReplace != nil && userInfoReplace.SnowflakeId != snowflakeId {

		userInfoReplace.Nick = userInfo.Nick
		userInfoReplace.Avatar = userInfo.Avatar

		err = data.UserInfoReplace(tx, userInfoReplace, userInfo.SnowflakeId)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("替换用户信息失败: %v", err)})
		}

		// 替换导入的用户信息

		err := data.UpdateIncomeExpense(tx, userInfoReplace.SnowflakeId, userInfo.SnowflakeId)

		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("替换用户收支明细信息失败: %v", err)})
		}

		// 删除原来的用户信息
		err = data.DeleteUser(tx, snowflakeId)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("删除用户信息失败: %v", err)})
		}
	} else {
		err = data.UpdateUserBySnowflakeId(tx, userInfo)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("修改用户信息失败: %v", err)})
		}
	}

	data.Commit(tx)
	return c.JSON(tools.Response{
		Code:    0,
		Message: "成功",
		Result:  struct{}{},
	})
}
