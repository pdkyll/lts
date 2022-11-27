package controller

import (
	"context"
	"strconv"
	
	"github.com/gofiber/fiber/v2"

	"gitlab.com/m0ta/lts/app/model"
	"gitlab.com/m0ta/lts/app/service"
	//"gitlab.com/m0ta/lts/app/utils"
)

// TokenController ...
type TokenController struct {
	ctx      context.Context
	services *service.Manager
	//logger   *logger.Logger
}

// NewTokens creates a new token controller.
func NewTokens(ctx context.Context, services *service.Manager) *TokenController {
	return &TokenController{
		ctx:      ctx,
		services: services,
		//logger:   logger,
	}
}

// Create creates new token
func (ctr *TokenController) Create(ctx *fiber.Ctx) error {
	token := &model.Token{}
	if err := ctx.BodyParser(token); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Check validation
	err := ctr.services.Token.Validate(ctx.Context(), token)
	if err != nil {
		return ctx.Status(412).JSON(fiber.Map{"status": "error", "message": "Incorrect data", "data": err})
	}
	
	// Create token
	token, err = ctr.services.Token.Create(ctx.Context(), token)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create token", "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Token created", "data": token})
}

// Get returns token by ID
func (ctr *TokenController) Get(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 0, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't parse Token ID", "data": nil})
	}

	token, err := ctr.services.Token.Get(ctx.Context(), id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "Token found", "data": token})
}

// Update token by token JSON
func (ctr *TokenController) Update(ctx *fiber.Ctx) error {
	token := &model.Token{}
	if err := ctx.BodyParser(token); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	if !(token.ID > 0) {
		id, err := strconv.ParseUint(ctx.Params("id"), 0, 32)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't parse Token ID", "data": nil})
		}
		token.ID = id
	}

	// Check validation
	err := ctr.services.Token.Validate(ctx.Context(), token)
	if err != nil {
		return ctx.Status(412).JSON(fiber.Map{"status": "error", "message": "Incorrect data", "data": err})
	}

	// Update token
	updToken, err := ctr.services.Token.Update(ctx.Context(), token)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't update token", "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Token updated", "data": updToken})
}

// Delete deletes token by ID
func (ctr *TokenController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 0, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't parse Token ID", "data": nil})
	}

	_, err = ctr.services.Token.Get(ctx.Context(), id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}

	err = ctr.services.Token.Delete(ctx.Context(), id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't delete token", "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "Token deleted", "data": nil})
}

// List returns token list
func (ctr *TokenController) List(ctx *fiber.Ctx) error {
	// _, err := uuid.Parse(ctx.Locals("user_id").(string))
	// if err != nil {
	// 	return ctx.Status(403).JSON(fiber.Map{"status": "error", "message": "Couldn't parse user UUID", "data": nil})
	// }

	tokens, err := ctr.services.Token.List(ctx.Context())
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't list tokens", "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "Token list found", "data": tokens})
}

// Verify verifies token
func (ctr *TokenController) Verify(ctx *fiber.Ctx) error {
	// type DomainToken struct {
	// 	Domain 	string
	// 	Token	string
	// }
	// input := &DomainToken{}
	// if err := ctx.BodyParser(input); err != nil {
	// 	return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	// }
	domain := ctx.Params("domain")

	token, err := ctr.services.Token.GetByDomain(ctx.Context(), domain)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}

	result, err := ctr.services.Token.Verify(ctx.Context(), token)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Token verified", "data": result})
}