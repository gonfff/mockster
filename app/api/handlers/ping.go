package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PingHandler is the handler for the /ping route
type PingHandler struct {
	e *echo.Echo
}

// NewPingHandler creates a new PingHandler
func NewPingHandler(e *echo.Echo) *PingHandler {
	return &PingHandler{e: e}
}

// RegisterRoutes registers the routes for the handler
func (h *PingHandler) RegisterRoutes() {
	h.e.GET("/ping", h.Ping)
}

// Ping is the handler to check if the application is running and healthy
func (h *PingHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
