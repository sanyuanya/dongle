package pc

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/tools"
)

func OrderCallback(c fiber.Ctx) error {

	taskId := c.FormValue("taskId")
	sign := c.FormValue("sign")
	param := c.FormValue("param")


	

	return c.JSON(tools.Response{
		Code:    0,
		Message: taskId,
		Result:  map[string]any{},
	})
}
