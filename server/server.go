package server

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v3"
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
	svr.argumentParse()

	svr.app.Get("/api/metrics", svr.handleMetrics)
	svr.app.Listen(":9227")
}

func (svr *server) argumentParse() {
	port := flag.String("port", "9227", "htpp server port ex) -port=8080")
	flag.Parse()

	svr.config.Port = *port
}
