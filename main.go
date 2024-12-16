package main

import (
	"be-kreditkita/src/config"

	"github.com/gofiber/fiber/v2"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	config.Connect()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World Yohanes!")
	})

	app.Listen(":8080")
}
