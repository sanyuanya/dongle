package pc

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/pay"
)

func Certificates(c fiber.Ctx) error {
	pay.Certificates()
	return c.JSON("")
}
