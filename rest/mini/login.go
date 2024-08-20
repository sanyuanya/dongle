package mini

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
	"github.com/sanyuanya/dongle/wechat"
)

func MiniLogin(c fiber.Ctx) error {

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

	miniLoginRequest := new(entity.MiniLoginRequest)

	err := c.Bind().Body(miniLoginRequest)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	code2SessionResp, err := wechat.Code2Session(miniLoginRequest.JsCode)
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取openid失败: %v", err)})
	}

	if code2SessionResp.ErrCode != 0 {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取openid失败: %v", code2SessionResp.ErrMsg)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	snowflakeId, err := data.FindOpenId(tx, code2SessionResp.OpenID)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("openid查询失败: %v", err)})
	}

	if snowflakeId == "" {
		snowflakeId = tools.SnowflakeUseCase.NextVal()

		registerUserRequest := &entity.RegisterUserRequest{
			SnowflakeId: snowflakeId,
			OpenId:      code2SessionResp.OpenID,
			SessionKey:  code2SessionResp.SessionKey,
		}

		err := data.RegisterUser(tx, registerUserRequest)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("openid注册失败: %v", err)})
		}
	} else {
		err := data.UpdateSessionKey(tx, code2SessionResp.OpenID, code2SessionResp.SessionKey)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("openid更新失败: %v", err)})
		}
	}

	token, err := tools.GenerateToken(snowflakeId, "user")
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50004, Message: fmt.Sprintf("生成token失败: %v", err)})
	}

	// 保存一下 token 方便测试
	err = data.UpdateUserApiToken(tx, snowflakeId, token)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50005, Message: fmt.Sprintf("更新token失败: %v", err)})
	}

	c.Response().Header.Set("Authorization", token)

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "成功",
		Result: map[string]any{
			"openid":       code2SessionResp.OpenID,
			"session_key":  code2SessionResp.SessionKey,
			"snowflake_id": snowflakeId,
		},
	})
}
