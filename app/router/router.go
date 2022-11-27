package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"go.opentelemetry.io/otel/attribute"

	"gitlab.com/m0ta/lts/app/controller"
	"gitlab.com/m0ta/lts/app/middleware"
	"gitlab.com/m0ta/lts/app/log"
)

// SetupWelcomeRoutes setup router api
func SetupWelcomeRoutes(ctx context.Context, app *fiber.App, c *controller.WelcomeController) fiber.Router{
	_, traceAPI := log.Tracer.Start(ctx, "api-middleware")
	defer traceAPI.End()
	app.Use(func(c *fiber.Ctx) error {
		//log := c.Hostname()
		log.Logger.Info("fiber",
			//zap.String("ctx", c.String()),
			//zap.String("str", c.IPs()[0]),
			zap.Int("status", c.Response().Header.StatusCode()),
			zap.ByteString("method", c.Request().Header.Method()),
			zap.String("path", c.Path()),
			zap.Strings("ips", c.IPs()),
			zap.String("requestHeader", c.Request().Header.String()),
			zap.ByteString("requestBody", c.Request().Body()),
			zap.String("responseHeader", c.Response().Header.String()),
			zap.ByteString("responseBody", c.Response().Body()),
		)

		traceAPI.SetAttributes(attribute.StringSlice("ips", c.IPs()))
		traceAPI.SetAttributes(attribute.Int("status", c.Response().Header.StatusCode()))
		traceAPI.SetAttributes(attribute.String("path", c.Path()))
		traceAPI.SetAttributes(attribute.String("method", string(c.Request().Header.Method())))

		return c.Next()
   	})
	api := app.Group("/api", logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05",
		Format: "[${time}] ${status} - ${latency} ${ips} ${method} ${path}\n",
		//Format: "[${time}] ${pid} ${status} - ${method} ${path} (${latency})\n",
		//Output:	os.Stderr,
	}))
	//api := app.Group("/api")
	api.Get("/", c.Welcome)

	return api
}

// SetupUserRoutes setup router api
func SetupUserRoutes(api fiber.Router, c *controller.UserController) {

	auth := api.Group("/")
	auth.Post("/signup", 	c.SignUp)
	auth.Post("/signin", 	c.SignIn)
	
	router := api.Group("/user").Use(middleware.Auth)
	router.Get("/", c.List)
	router.Get("/:id", c.Get)
	router.Patch("/", c.Update)
	router.Delete("/", c.Delete)
	router.Patch("/change-password", c.ChangePassword)
}

// SetupTokenRoutes setup router api
func SetupTokenRoutes(api fiber.Router, c *controller.TokenController) {

	router := api.Group("/token")
	//router := api.Group("/token").Use(middleware.Auth)
	router.Get("/", c.List)
	router.Get("/:id", c.Get)
	router.Post("/", c.Create)
	router.Patch("/:id", c.Update)
	router.Delete("/:id", c.Delete)

	router.Get("/verify/:domain", c.Verify)
	//router.Post("/verify", c.Verify)
}