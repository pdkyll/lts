package controller

import (
	"context"
	
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"gitlab.com/m0ta/lts/app/model"
	"gitlab.com/m0ta/lts/app/service"
	"gitlab.com/m0ta/lts/app/utils"
)

// UserController ...
type UserController struct {
	ctx      context.Context
	services *service.Manager
	//logger   *logger.Logger
}

// NewUsers creates a new user controller.
func NewUsers(ctx context.Context, services *service.Manager) *UserController {
	return &UserController{
		ctx:      ctx,
		services: services,
		//logger:   logger,
	}
}

// Create creates new user
func (ctr *UserController) Create(ctx *fiber.Ctx) error {
	user := &model.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Check validation
	err := ctr.services.User.Validate(ctx.Context(), user)
	if err != nil {
		return ctx.Status(412).JSON(fiber.Map{"status": "error", "message": "Incorrect data", "data": err})
	}
	
	// Create user
	user, err = ctr.services.User.Create(ctx.Context(), user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User created", "data": user})
}

// Get returns user by ID
func (ctr *UserController) Get(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(403).JSON(fiber.Map{"status": "error", "message": "Couldn't parse user UUID", "data": nil})
	}

	user, err := ctr.services.User.Get(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// Update user by user JSON
func (ctr *UserController) Update(ctx *fiber.Ctx) error {
	user := &model.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	userID, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(403).JSON(fiber.Map{"status": "error", "message": "Couldn't parse user UUID", "data": nil})
	}

	// Check validation
	// err = ctr.services.User.Validate(ctx.Context(), user)
	// if err != nil {
	// 	return ctx.Status(412).JSON(fiber.Map{"status": "error", "message": "Incorrect data", "data": err})
	// }

	// Update user
	user.ID = userID
	updUser, err := ctr.services.User.Update(ctx.Context(), user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't update user", "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User updated", "data": updUser})
}

// Delete deletes user by ID
func (ctr *UserController) Delete(ctx *fiber.Ctx) error {
	user := &model.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	userID, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(403).JSON(fiber.Map{"status": "error", "message": "Couldn't parse user UUID", "data": nil})
	}

	foundUser, err := ctr.services.User.Get(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}

	// Verify password
	err = utils.VerifyPassword(foundUser.EncryptedPassword, user.Password)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{"status": "error", "message": "Incorrect email or password", "data": err})
	}

	err = ctr.services.User.Delete(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't delete user", "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "User deleted", "data": nil})
}

// List returns user list
func (ctr *UserController) List(ctx *fiber.Ctx) error {
	_, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(403).JSON(fiber.Map{"status": "error", "message": "Couldn't parse user UUID", "data": nil})
	}

	users, err := ctr.services.User.List(ctx.Context())
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't list users", "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "User list found", "data": users})
}

// SignIn ...
func (ctr *UserController) SignIn(ctx *fiber.Ctx) error {
	user := &model.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	foundUser, err := ctr.services.User.GetByEmail(ctx.Context(), user.Email)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}

	// Verify password
	err = utils.VerifyPassword(foundUser.EncryptedPassword, user.Password)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{"status": "error", "message": "Incorrect email or password", "data": err})
	}

	// Get user token
	token, err := foundUser.GenerateToken()//utils.TokenGenerate(foundUser)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't generate token", "data": err})
	}
	foundUser.Token = token

	return ctx.JSON(fiber.Map{"status": "success", "message": "User logined", "data": foundUser})
}

// SignUp ...
func (ctr *UserController) SignUp(ctx *fiber.Ctx) error {
	user := &model.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Check validation
	err := ctr.services.User.Validate(ctx.Context(), user)
	if err != nil {
		return ctx.Status(412).JSON(fiber.Map{"status": "error", "message": "Incorrect data", "data": err})
	}

	// Create user
	updUser, err := ctr.services.User.Create(ctx.Context(), user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	// Get user token
	token, err := updUser.GenerateToken()//utils.TokenGenerate(updUser)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't generate token"), "data": err})
	}
	updUser.Token = token

	return ctx.JSON(fiber.Map{"status": "success", "message": "User created", "data": updUser})
}

// ChangePassword ...
func (ctr *UserController) ChangePassword(ctx *fiber.Ctx) error {
	user := &model.User{}
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	// Check validation
	err := ctr.services.User.Validate(ctx.Context(), user)
	if err != nil {
		return ctx.Status(412).JSON(fiber.Map{"status": "error", "message": "Incorrect data", "data": err})
	}

	userID, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(403).JSON(fiber.Map{"status": "error", "message": "Couldn't parse user UUID", "data": nil})
	}

	// Find user
	updUser, err := ctr.services.User.Get(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "Data couldn't found", "data": nil})
	}

	// Change user password
	updUser, err = ctr.services.User.ChangePassword(ctx.Context(), updUser, user.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't change user password", "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User password changed", "data": updUser})
}