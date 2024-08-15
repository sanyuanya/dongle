package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/sanyuanya/dongle/rest/mini"
	"github.com/sanyuanya/dongle/rest/pc"
	"github.com/sanyuanya/dongle/tools"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()
	app.Use(recover.New())

	app.Use(func(c fiber.Ctx) error {

		var jsonBody map[string]interface{}

		// Only read and log the body for non-GET requests
		if c.Method() == fiber.MethodPost {

			contentType := c.Get("Content-Type")

			if contentType == "application/json" {
				// Read the request body
				body, err := io.ReadAll(bytes.NewReader(c.Body()))
				if err != nil {
					return c.JSON(tools.Response{
						Code:    50010,
						Message: "Error reading body",
						Result:  struct{}{},
					})
				}

				// Check if the body is empty
				if len(body) != 0 {
					// Parse the request body as JSON
					if err := json.Unmarshal(body, &jsonBody); err != nil {
						return c.JSON(tools.Response{
							Code:    50011,
							Message: "Error unmarshalling JSON",
							Result:  struct{}{},
						})
					}
				}
				// Set the body back to the context
				c.Context().SetBody(body)
			}
		}

		tools.Logger.Info("Request:",
			"method", c.Method(),
			"path", c.OriginalURL(),
			"headers", c.GetReqHeaders(),
			"body", jsonBody,
		)

		// Continue to the next middleware/handler
		return c.Next()
	})
	// Define a route for the GET method on the root path '/'
	// pc
	app.Post("/api/pc/import", pc.ExcelImport)

	app.Post("/api/pc/login", pc.PcLogin)

	app.Get("/api/pc/userList", pc.UserList)

	app.Post("/api/pc/addWhite", pc.SetUpWhite)

	app.Get("/api/pc/withdrawalList", pc.WithdrawalList)

	app.Post("/api/pc/approvalWithdrawal", pc.ApprovalWithdrawal)

	app.Post("/api/pc/updateUserInfo", pc.UpdateUserInfo)

	app.Get("/api/pc/incomeList", pc.GetIncomeList)

	app.Get("/api/pc/batch/:batchId", pc.OutBatchNo)

	app.Get("/api/pc/batch/:batchId/transfer/:transferId", pc.OutTransferNo)

	app.Post("/api/pc/pay", pc.Pay)

	//mini program

	app.Post("/api/mini/login", mini.MiniLogin)

	app.Post("/api/mini/updateUserInfo", mini.UpdateUserInfo)

	app.Get("/api/mini/getUserInfo", mini.GetUserInfo)

	app.Post("/api/mini/setUserInfo", mini.SetUserInfo)

	app.Post("/api/mini/applyForWithdrawal", mini.ApplyForWithdrawal)

	app.Get("/api/mini/getWithdrawalList", mini.GetWithdrawalList)

	app.Get("/api/mini/incomeList", mini.GetIncomeList)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
