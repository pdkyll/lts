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
//func NewUsers(ctx context.Context, services *service.Manager, logger *logger.Logger) *UserController {
func NewUsers(ctx context.Context, services *service.Manager) *UserController {
	return &UserController{
		ctx:      ctx,
		services: services,
		//logger:   logger,
	}
}


// Create creates new user
func (ctr *UserController) Create(ctx *fiber.Ctx) error {

	user := *new(model.User)
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Review your input"), "data": err})

	}

	// check validation
	
	// hash password
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't hash password"), "data": err})

	}
	user.EncryptedPassword = hash
	
	// create user
	createdUser, err := ctr.services.User.Create(ctx.Context(), &user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't create user"), "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User created", "data": createdUser})
}

// Get returns user by ID
func (ctr *UserController) Get(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "could not parse user UUID"), "data": nil})
	}

	user, err := ctr.services.User.Get(ctx.Context(), userID)

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "could not get user"), "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// Update user by user JSON
func (ctr *UserController) Update(ctx *fiber.Ctx) error {
	user := *new(model.User)
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Review your input (UpdateUser)"), "data": err})
	}

	userID, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "could not parse user UUID"), "data": nil})
	}

	// check validation
	
	user.ID = userID
	
	// update user
	updatedUser, err := ctr.services.User.Update(ctx.Context(), &user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't update user"), "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User updated", "data": updatedUser})
}

// Delete deletes user by ID
func (ctr *UserController) Delete(ctx *fiber.Ctx) error {
	m := *new(model.User)
	if err := ctx.BodyParser(&m); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Review your input (UpdateUser)"), "data": err})
	}

	userID, err := uuid.Parse(ctx.Locals("user_id").(string))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "could not parse user UUID"), "data": nil})
	}

	user, err := ctr.services.User.Get(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "no user found"), "data": err})
	}

	// verify password
	err = utils.VerifyPassword(user.EncryptedPassword, m.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "incorrect email or password"), "data": err})
	}

	err = ctr.services.User.Delete(ctx.Context(), userID)

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "could not delete user"), "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "User deleted", "data": nil})
}

// List returns user by ID
func (ctr *UserController) List(ctx *fiber.Ctx) error {
	users, err := ctr.services.User.List(ctx.Context())

	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "could not get user"), "data": nil})
	}
	
	return ctx.JSON(fiber.Map{"status": "success", "message": "User's list found", "data": users})
}

// SignIn ...
func (ctr *UserController) SignIn(ctx *fiber.Ctx) error {
	type LoginInput struct {
		Email 		string `json:"email"`
		Password 	string `json:"password"`
	}

	var input *LoginInput

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Review your input"), "data": err})
	}

	// find user by email
	foundUser, err := ctr.services.User.GetUserByEmail(ctx.Context(), input.Email)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't find user by email"), "data": err})
	}

	// verify password
	err = utils.VerifyPassword(foundUser.EncryptedPassword, input.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Incorrect email or password"), "data": err})
	}

	// get user token
	token, err := utils.TokenGenerate(foundUser)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't generate token"), "data": 530})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User signined", "data": foundUser, "token": token})
}

// SignUp ...
func (ctr *UserController) SignUp(ctx *fiber.Ctx) error {
	user := *new(model.User)
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Review your input"), "data": err})

	}

	// check validation
	
	// hash password
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't hash password"), "data": err})

	}
	user.EncryptedPassword = hash
	
	// создаем пользователя
	createdUser, err := ctr.services.User.Create(ctx.Context(), &user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't create user"), "data": err})
	}

	// get user token
	token, err := utils.TokenGenerate(createdUser)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": utils.ErrorString(err, "Couldn't generate token"), "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User created", "data": createdUser, "token": token})
}