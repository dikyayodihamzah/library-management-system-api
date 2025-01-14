package controller

import (
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *controller) signUp(ctx *fiber.Ctx) error {
	api := new(user.SignUpRequest)
	if err := lib.BodyParser(ctx, api); err != nil {
		return exception.Handler(ctx, err)
	}

	res, err := c.UserService.SignUp(ctx.Context(), api)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.Created(ctx, res)
}

func (c *controller) login(ctx *fiber.Ctx) error {
	api := new(user.LoginRequest)
	if err := lib.BodyParser(ctx, api); err != nil {
		return exception.Handler(ctx, err)
	}

	res, err := c.UserService.Login(ctx.Context(), api)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HTTPOnly: true,
	})

	return lib.OK(ctx, res)
}

func (c *controller) logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
	})

	return lib.OK(ctx)
}

func (c *controller) findAllUsers(ctx *fiber.Ctx) error {
	filter := new(model.QueryParam)
	if err := ctx.QueryParser(filter); err != nil {
		return exception.Handler(ctx, err)
	}

	res, total, err := c.UserService.FindAll(ctx.Context(), filter)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.Page(ctx, total, res)
}

func (c *controller) assignAdmin(ctx *fiber.Ctx) error {
	id := new(uuid.UUID)
	if res := lib.StrToUUID(ctx.Params("id")); res == nil || *res == uuid.Nil {
		return exception.Handler(ctx, exception.ErrorBadRequest("Invalid ID"))
	} else {
		id = res
	}

	res, err := c.UserService.AssignAdmin(ctx.Context(), *id)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.OK(ctx, res)
}
