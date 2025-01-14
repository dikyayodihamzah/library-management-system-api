package middleware

import (
	"strings"

	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	claims, err := lib.ParseJwt(lib.GetToken(c))
	if err != nil || claims == nil {
		res := model.Response{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		}

		if strings.Contains(err.Error(), "token is expired") {
			res.Message = "Token Expired"
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}

		return c.Status(fiber.StatusUnauthorized).JSON(res)
	}

	c.Locals("claims", claims)

	return c.Next()
}

func IsAdmin(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*lib.Claims)
	if !claims.IsAdmin {
		res := model.Response{
			Code:    fiber.StatusForbidden,
			Message: "Forbidden",
		}

		return c.Status(fiber.StatusUnauthorized).JSON(res)
	}

	return c.Next()
}
