package main

import (
	"authentification/controllers"
	"authentification/db"
	"authentification/repository"
	"authentification/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func init() {
	err := godotenv.Load("database.env")

	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	conn := db.NewConnection()
	defer conn.Close()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Hello World"})
	})

	usersRepo := repository.NewUsersRepository(conn)
	authController := controllers.NewAuthController(usersRepo)
	authRoutes := routes.NewAuthRoutes(authController)
	authRoutes.Install(app)

	log.Fatal(app.Listen(":8080"))
}
