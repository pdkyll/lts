package router

import (
	"github.com/gofiber/fiber/v2"

	"gitlab.com/m0ta/lts/app/controller"
	"gitlab.com/m0ta/lts/app/handler"
	"gitlab.com/m0ta/lts/app/middleware"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App, c *controller.UserController) fiber.Router{
	api := app.Group("/api")
	api.Get("/", handler.Hello)

	// Auth
	auth := app.Group("/auth")
	auth.Post("/signup", 	c.SignUp)
	auth.Post("/signin", 	c.SignIn)

	return api
}

// SetupRoutesForUser setup router api
func SetupRoutesForUser(api fiber.Router, c *controller.UserController) {
	
	user := api.Group("/user").Use(middleware.Auth)
	//user.Post("/", c.Create)
	user.Get("/", c.List)//.Use(middleware.Auth)
	user.Patch("/", c.Update)//.Use(middleware.Auth)
	user.Delete("/", c.Delete)//.Use(middleware.Auth)

	// users := api.Group("/users").Use(middleware.Auth)
	// users.Get("/", c.List)
}