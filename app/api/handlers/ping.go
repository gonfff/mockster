package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// type PingResponse struct {

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
