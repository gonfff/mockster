package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Ping is the handler to check if the application is running and healthy
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
