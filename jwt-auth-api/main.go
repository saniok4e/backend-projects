package main

import (

    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"jwt-auth-api/routes"
	"jwt-auth-api/database"
)

func main() {
	database.DBconn()

    app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8000",
        AllowCredentials: true, 
    }))

    routes.Setup(app) 

    app.Listen(":8000")
}