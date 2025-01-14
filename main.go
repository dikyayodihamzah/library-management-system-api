package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dikyayodihamzah/library-management-api/app/controller"
	"github.com/dikyayodihamzah/library-management-api/app/repository/bookrepo"
	"github.com/dikyayodihamzah/library-management-api/app/repository/borrowrepo"
	"github.com/dikyayodihamzah/library-management-api/app/repository/userrepo"
	"github.com/dikyayodihamzah/library-management-api/app/service/booksvc"
	"github.com/dikyayodihamzah/library-management-api/app/service/borrowsvc"
	"github.com/dikyayodihamzah/library-management-api/app/service/usersvc"
	"github.com/dikyayodihamzah/library-management-api/pkg/config/dbconfig"
	"github.com/dikyayodihamzah/library-management-api/pkg/lib"
	"github.com/dikyayodihamzah/library-management-api/pkg/transaction"
	"github.com/dikyayodihamzah/library-management-api/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	logger = utils.NewLogger()
)

// ========== ROUTES ==========
func listenRoutes(
	ctrl controller.Controller,
) {
	// setup fiber app
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			lib.PrintStackTrace(e)
		},
	}))

	// provide health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return lib.OK(c, utils.GetString("APP_NAME"))
	})

	// setup routes
	ctrl.Routes(app)

	// start server
	if err := app.Listen(fmt.Sprintf(":%s", utils.GetString("SERVER_PORT"))); err != nil {
		log.Fatalln("Failed to start server", err)
	}
}

func main() {
	time.Local = time.UTC

	// Initialization
	postgreDB := dbconfig.NewPool()
	txManager := transaction.New(postgreDB)

	// repository
	userRepository := userrepo.New(logger, postgreDB)
	bookRepository := bookrepo.New(logger, postgreDB)
	borrowRepository := borrowrepo.New(logger, postgreDB)

	// service
	validate := validator.New()
	userService := usersvc.New(logger, validate, txManager, userRepository)
	bookService := booksvc.New(validate, txManager, bookRepository)
	borrowService := borrowsvc.New(logger, validate, txManager, userRepository, bookRepository, borrowRepository)

	// controller
	ctrl := controller.New(userService, bookService, borrowService)

	// listen to routes
	listenRoutes(ctrl)
}
