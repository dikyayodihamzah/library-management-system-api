package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
	"github.com/gofiber/fiber/v2"
)

func (c *controller) borrowBook(ctx *fiber.Ctx) error {
	req := new(book.BorrowRequest)
	if err := ctx.BodyParser(req); err != nil {
		return exception.Handler(ctx, exception.ErrorBadRequest(err.Error()))
	}

	claims := ctx.Locals("claims").(*lib.Claims)
	req.UserID = *lib.StrToUUID(claims.Issuer)

	res, err := c.BorrowService.Borrow(ctx.Context(), req)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.Created(ctx, res)
}

func (c *controller) returnBook(ctx *fiber.Ctx) error {
	req := new(book.BorrowRequest)
	if err := ctx.BodyParser(req); err != nil {
		return exception.Handler(ctx, exception.ErrorBadRequest(err.Error()))
	}

	claims := ctx.Locals("claims").(*lib.Claims)
	req.UserID = *lib.StrToUUID(claims.Issuer)

	if err := c.BorrowService.Return(ctx.Context(), req); err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.OK(ctx)
}

func (c *controller) findAllBorrows(ctx *fiber.Ctx) error {
	filter := new(book.BorrowQuery)
	if err := ctx.QueryParser(filter); err != nil {
		return exception.Handler(ctx, exception.ErrorBadRequest(err.Error()))
	}

	claims := ctx.Locals("claims").(*lib.Claims)
	if !claims.IsAdmin {
		filter.UserID = *lib.StrToUUID(claims.Issuer)
	}

	res, total, err := c.BorrowService.FindAll(ctx.Context(), filter)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.Page(ctx, total, res)
}

func (c *controller) generateBorrowExcel(ctx *fiber.Ctx) error {
	filter := new(book.BorrowQuery)
	if err := ctx.QueryParser(filter); err != nil {
		return exception.Handler(ctx, exception.ErrorBadRequest(err.Error()))
	}

	timezone, _ := strconv.Atoi(ctx.Query("timezone"))

	file, err := c.BorrowService.GenerateExcel(ctx.Context(), filter, timezone)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	name := "excel-" + time.Now().Format("20060102150405") + ".xlsx"
	ctx.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", name))
	ctx.Set(fiber.HeaderContentLength, fmt.Sprint(file.Len()))
	return ctx.SendStream(file)
}
