package routes

import (
	"be-kreditkita/src/controllers"
	"be-kreditkita/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/login", controllers.Login)
	api.Post("/consumers", controllers.CreateConsumer)

	app.Use(middlewares.IsAuth)

	api.Get("/user", controllers.User)
	api.Post("/logout", controllers.Logout)

	api.Get("/consumers", controllers.AllConsumers)
	api.Get("/consumer/:id", controllers.GetConsumer)
	api.Put("/consumer/:id", controllers.UpdateConsumer)
	api.Delete("/consumer/:id", controllers.DeleteConsumer)

	api.Get("/limits", controllers.AllLimit)
	api.Get("/limits/:id", controllers.GetLimit)
	api.Put("/limits/:id", controllers.UpdateLimit)
	api.Delete("/limits/:id", controllers.DeleteLimit)
	api.Put("/limits/consumer/:id", controllers.UpdateLimitConsumer)
	api.Get("/limits/consumer/:id", controllers.GetLimitByConsumerId)

	app.Get("/transactions", controllers.AllTransactions)
	app.Get("/transaction/:id", controllers.GetTransactions)
	app.Post("/transactions", controllers.CreateTransactions)
	app.Put("/transaction/:id", controllers.UpdateTransactions)
	app.Delete("transaction/:id", controllers.DeleteTransactions)
	app.Get("/transactions/consumer/:id", controllers.GetTransactionsByConsumerId)
}
