package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PingHandler is the handler for the /ping route
type PingHandler struct{}

// NewPingHandler creates a new PingHandler
func NewPingHandler() *MockHandler {
	return &MockHandler{}
}

// RegisterRoutes registers the routes for the handler
func (h *PingHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/ping", h.Ping)
}

// Ping is the handler to check if the application is running and healthy
func (h *PingHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
