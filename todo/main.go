package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("Hello, world👋👋!")
	})

	log.Fatal(app.Listen(":3000"))
}
