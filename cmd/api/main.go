package main

import (
	"context"
	//"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"

	"gitlab.com/m0ta/lts/app/config"
	"gitlab.com/m0ta/lts/app/controller"
	"gitlab.com/m0ta/lts/app/router"
	"gitlab.com/m0ta/lts/app/service"
	"gitlab.com/m0ta/lts/app/store"
	"gitlab.com/m0ta/lts/app/utils"
	"gitlab.com/m0ta/lts/app/log"
)

func main() {
	cfg 	:= config.Get()
	//logger 	:= logger.New(cfg)//newLogger(cfg)
	log.New(cfg)
    defer log.Logger.Sync()
	
	if err := run(cfg, log.Logger); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

// func newLogger(cfg *config.Config) *zap.Logger {
// 	pe 		:= zap.NewProductionEncoderConfig()
// 	level 	:= zapcore.Level(cfg.LogLevel)
// 	//pe := zap.NewDevelopmentEncoderConfig()
//     //fileEncoder := zapcore.NewJSONEncoder(pe)
	
//     pe.EncodeTime 	= zapcore.ISO8601TimeEncoder
// 	pe.EncodeLevel 	= zapcore.CapitalColorLevelEncoder
// 	pe.EncodeCaller = zapcore.FullCallerEncoder
	
//     consoleEncoder 	:= zapcore.NewConsoleEncoder(pe)

//  	core := zapcore.NewTee(
//         //zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zap.DebugLevel),
//         zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
//     )
// 	return zap.New(core, zap.AddCaller())
// }

func run(cfg *config.Config, logger *zap.Logger) error {
	ctx := context.Background()
	// Start uptrace
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		//uptrace.WithDSN("https://<key>@uptrace.dev/<project_id>"),
		//uptrace.WithServiceName("<service_name>"),
		//uptrace.WithServiceVersion("<version>"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	// Start root span of otel
	//tracer := otel.Tracer("lts")
	ctx, rootSpan := log.Tracer.Start(ctx, "main")
	defer rootSpan.End()

	// Init repository store (with postgresql inside)
	store, err := store.New(ctx, logger)
	if err != nil {
		return utils.ErrorWrap(err, "store.New failed")
	}
	
	// Init service manager
	serviceManager, err := service.NewManager(ctx, store, logger)
	if err != nil {
		return utils.ErrorWrap(err, "manager.New failed")
	}

	// Init controllers
	cWelcome	:= controller.NewWelcome(ctx, serviceManager)
	cUser 		:= controller.NewUsers(ctx, serviceManager)
	cToken 		:= controller.NewTokens(ctx, serviceManager)
	
	// Initialize Fiber instance
	app := fiber.New()
	app.Use(cors.New())

	api := router.SetupWelcomeRoutes(ctx, app, cWelcome)
	router.SetupUserRoutes(api, cUser)
	router.SetupTokenRoutes(api, cToken)

	s := app.Stack()
	for _, v := range s {
		for _, w := range v {
			if (w.Method == "GET") || (w.Method == "POST") || (w.Method == "DELETE") || (w.Method == "PATCH") {
				serviceManager.Logger.Debug(
					"available api methods:",
					zap.String("Method", w.Method),
					zap.String("Path", w.Path),
				)
			}
		}
	}

	// Start api server
	logger.Fatal(app.Listen(cfg.APIAddr).Error())

	return nil
}