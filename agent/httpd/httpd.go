package httpd

import (
	"github.com/gofiber/fiber/v3"
)

type app struct {
	*fiber.App
}

func New() *app {
	return &app{fiber.New()}
}

func (a *app) Close() {
	a.Shutdown()
}

func (a *app) Serve() {
	a.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	a.Listen(":9227")
}
