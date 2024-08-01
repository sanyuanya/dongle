package main

import "github.com/gofiber/fiber/v3"

func GetUserInfo(c fiber.Ctx) error {

	return c.JSON(map[string]string{})
}
