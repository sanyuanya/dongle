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
			c.JSON(tools.Response{
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

	fmt.Printf("snowflakeId: %v\n", snowflakeId)
	if snowflakeId == 0 {

		fmt.Printf("snowflakeId 2: %v\n", snowflakeId)

		snowflakeId = tools.SnowflakeUseCase.NextVal()

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

		fmt.Printf("snowflakeId 3: %v\n", snowflakeId)

		err := data.UpdateSessionKey(code2SessionResp.OpenID, code2SessionResp.SessionKey)
		if err != nil {
			panic(fmt.Errorf("openid更新失败: %v", err))
		}
	}

	token, err := tools.GenerateToken(snowflakeId, "user")

	if err != nil {
		panic(fmt.Errorf("生成token失败: %v", err))
	}

	// 保存一下 token 方便测试

	err = data.UpdateUserApiToken(snowflakeId, token)

	if err != nil {
		panic(fmt.Errorf("更新token失败: %v", err))
	}

	c.Response().Header.Set("Authorization", token)

	fmt.Printf("snowflakeId 4: %v\n", snowflakeId)
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
