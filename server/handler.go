package server

import (
	"helloServer/agent/measure"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Request struct {
}

type Response struct {
	Success bool             `json:"success"`
	Reason  string           `json:"reason"`
	Data    *measure.Measure `json:"data"`
	Elapse  string           `json:"elapse"`
}

func (svr *server) handleMetrics(c fiber.Ctx) error {
	tick := time.Now()
	rsp := Response{Success: false, Reason: "not specified"}

	val := measure.Get("l")

	rsp.Success = true
	rsp.Reason = "success"
	rsp.Data = val
	rsp.Elapse = time.Since(tick).String()
	return c.Status(200).JSON(rsp)
}
