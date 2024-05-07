package server

import (
	"fmt"
	"helloServer/measure"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Request struct {
	Period int64 `json:"period"`
}

type Response struct {
	Success bool        `json:"success"`
	Reason  string      `json:"reason"`
	Data    interface{} `json:"data,omitempty"`
	Elapse  string      `json:"elapse"`
}

func (svr *server) handleMetrics(c fiber.Ctx) error {
	now := time.Now()
	rsp := Response{Success: false, Reason: "not specified"}

	data := measure.Get("l")
	if data == nil {
		rsp.Reason = "data is empty (error : 'agent' or 'cache')"
		rsp.Elapse = time.Since(now).String()
		return c.Status(http.StatusInternalServerError).JSON(rsp)
	}

	metric := strings.ToLower(c.Params("metric"))
	switch metric {
	case "all":
		rsp.Data = data
	case "cpu":
		rsp.Data = data.Cpu
	case "memory":
		rsp.Data = data.Memory
	case "disk":
		rsp.Data = data.Disk
	case "network":
		rsp.Data = data.Network
	case "system":
		rsp.Data = data.System
	default:
		rsp.Reason = fmt.Sprintf("invalid metric: '%s'", metric)
		rsp.Elapse = time.Since(now).String()
		return c.Status(http.StatusBadRequest).JSON(rsp)
	}

	rsp.Success = true
	rsp.Reason = "success"
	rsp.Elapse = time.Since(now).String()
	return c.Status(http.StatusOK).JSON(rsp)
}

func (svr *server) handlePeriod(c fiber.Ctx) error {
	now := time.Now()
	req := Request{}
	rsp := Response{Success: false, Reason: "not specified"}

	if err := c.Bind().JSON(&req); err != nil {
		rsp.Reason = err.Error()
		rsp.Elapse = time.Since(now).String()
		return c.Status(http.StatusBadRequest).JSON(rsp)
	}

	if req.Period < 1 {
		rsp.Reason = "invalid period ( period > 0 )"
		rsp.Elapse = time.Since(now).String()
		return c.Status(http.StatusBadRequest).JSON(rsp)
	}

	if err := measure.Publish("period", time.Duration(req.Period)*time.Second); err != nil {
		rsp.Reason = err.Error()
		rsp.Elapse = time.Since(now).String()
		return c.Status(http.StatusInternalServerError).JSON(rsp)
	}

	rsp.Success = true
	rsp.Reason = "success"
	rsp.Elapse = time.Since(now).String()
	return c.Status(http.StatusOK).JSON(rsp)
}
