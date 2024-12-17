package main

import (
	"be-kreditkita/src/config"
	"be-kreditkita/src/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	config.Connect()
	helpers.Migrate()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World Yohanes!")
	})

	app.Listen(":8080")
}
