package router

import (
	"github.com/Fajar3108/mafi-course-be/config"
	errorhandler "github.com/Fajar3108/mafi-course-be/pkg/error-handler"
	"github.com/Fajar3108/mafi-course-be/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/spf13/viper"
)

func SetupRoutes() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.GlobalErrorHandler,
	})

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: viper.GetString(config.CookieSecretKey),
	}))

	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Welcome to the Shopee Clone API")
	})

	AuthRouter(v1.Group("/auth"))

	v1.Use(middleware.JWTMiddleware())
	CategoryRouter(v1)

	return app
}
