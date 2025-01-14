package controller

import (
	"github.com/dikyayodihamzah/library-management-api/app/service/booksvc"
	"github.com/dikyayodihamzah/library-management-api/app/service/borrowsvc"
	"github.com/dikyayodihamzah/library-management-api/app/service/usersvc"
	"github.com/dikyayodihamzah/library-management-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Routes(app *fiber.App)
}

type controller struct {
	UserService   usersvc.UserService
	BookService   booksvc.BookService
	BorrowService borrowsvc.BorrowService
}

func New(
	userService usersvc.UserService,
	bookService booksvc.BookService,
	borrowService borrowsvc.BorrowService,
) Controller {
	return &controller{
		UserService:   userService,
		BookService:   bookService,
		BorrowService: borrowService,
	}
}

func (c *controller) Routes(app *fiber.App) {
	app.Post(("/sign-up"), c.signUp)
	app.Post("/login", c.login)
	app.Post("/logout", c.logout)

	userAPI := app.Group("/users").Use(middleware.IsAuthenticated, middleware.IsAdmin)
	userAPI.Get("/", c.findAllUsers)
	userAPI.Post("/assign-admin/:id", c.assignAdmin)

	bookAPI := app.Group("/books").Use(middleware.IsAuthenticated)
	bookAPI.Post("/", middleware.IsAdmin, c.createBook)
	bookAPI.Get("/", c.findAllBooks)
	bookAPI.Get("/:id", c.findBookByID)
	bookAPI.Put("/:id", middleware.IsAdmin, c.updateBook)
	bookAPI.Delete("/:id", middleware.IsAdmin, c.deleteBook)

	borrowAPI := app.Group("/borrows").Use(middleware.IsAuthenticated)
	borrowAPI.Post("/", c.borrowBook)
	borrowAPI.Post("/return", c.returnBook)
	borrowAPI.Get("/", c.findAllBorrows)
	borrowAPI.Get("/excel", c.generateBorrowExcel)
}
