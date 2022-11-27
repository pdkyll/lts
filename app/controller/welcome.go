package controller

import (
	"context"
	
	"github.com/gofiber/fiber/v2"
	
	"gitlab.com/m0ta/lts/app/service"
	"gitlab.com/m0ta/lts/app/config"
)

// WelcomeController ...
type WelcomeController struct {
	ctx      context.Context
	//services *service.Manager
}

// NewWelcome creates a new welcome controller.
func NewWelcome(ctx context.Context, services *service.Manager) *WelcomeController {
	return &WelcomeController{
		ctx:      ctx,
		//services: services,
	}
}

// Welcome notifies the service is working
func (ctr *WelcomeController) Welcome(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "success", "message": "License Token Service operational!", "data": config.Version})
}