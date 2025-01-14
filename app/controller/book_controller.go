package controller

import (
	"github.com/dikyayodihamzah/library-management-api/pkg/exception"
	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/model"
	"github.com/dikyayodihamzah/library-management-api/pkg/model/web/book"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *controller) createBook(ctx *fiber.Ctx) error {
	api := new(book.BookRequest)
	if err := lib.BodyParser(ctx, api); err != nil {
		return exception.Handler(ctx, err)
	}

	res, err := c.BookService.Create(ctx.Context(), api)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.Created(ctx, res)
}

func (c *controller) findAllBooks(ctx *fiber.Ctx) error {
	filter := new(model.QueryParam)
	if err := ctx.QueryParser(filter); err != nil {
		return exception.Handler(ctx, err)
	}

	books, total, err := c.BookService.FindAll(ctx.Context(), filter)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.Page(ctx, total, books)
}

func (c *controller) findBookByID(ctx *fiber.Ctx) error {
	id := lib.StrToUUID(ctx.Params("id"))
	if id == nil || *id == uuid.Nil {
		return exception.ErrorBadRequest("Invalid ID")
	}

	book, err := c.BookService.FindByID(ctx.Context(), *id)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.OK(ctx, book)
}

func (c *controller) updateBook(ctx *fiber.Ctx) error {
	id := lib.StrToUUID(ctx.Params("id"))
	if id == nil || *id == uuid.Nil {
		return exception.ErrorBadRequest("Invalid ID")
	}

	api := new(book.BookRequest)
	if err := lib.BodyParser(ctx, api); err != nil {
		return exception.Handler(ctx, err)
	}

	res, err := c.BookService.Update(ctx.Context(), *id, api)
	if err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.OK(ctx, res)
}

func (c *controller) deleteBook(ctx *fiber.Ctx) error {
	id := lib.StrToUUID(ctx.Params("id"))
	if id == nil || *id == uuid.Nil {
		return exception.ErrorBadRequest("Invalid ID")
	}

	if err := c.BookService.Delete(ctx.Context(), *id); err != nil {
		return exception.Handler(ctx, err)
	}

	return lib.OK(ctx)
}
