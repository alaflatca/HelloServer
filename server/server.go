package server

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type Config struct {
	Port string `json:"port"`
}

type server struct {
	app    *fiber.App
	config Config
}

func New() *server {
	return &server{app: fiber.New()}
}

func (svr *server) Close() {
	log.Println("Fiber Shutdown")
	svr.app.Shutdown()
}

func (svr *server) Start() {
	svr.app.Use(recover.New())

	svr.app.Get("/api/metrics/:metric", svr.handleMetrics)
	svr.app.Put("/api/period", svr.handlePeriod)

	svr.app.Listen(":9227")
}
