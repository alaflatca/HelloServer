package server

import (
	"log"
	"net/http"

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

func (svr *server) Close() error {
	log.Println("Fiber Shutdown")
	return svr.app.Shutdown()
}

func (svr *server) Start() error {
	svr.registerRoutes()
	return svr.app.Listen(":9227")
}

func (svr *server) registerRoutes() {
	svr.app.Use(recover.New())

	svr.app.Get("/api/metrics", svr.handleMetrics)
	svr.app.Get("/api/metrics/:metric", svr.handleMetrics)
	svr.app.Get("/api/processes", svr.handleProcessList)
	svr.app.Put("/api/agent/period", svr.handlePeriod)

	svr.app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}
