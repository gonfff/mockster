package handlers

import (
	"io"
	"net/http"

	"github.com/gonfff/mockster/app/services"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// MockHandler handles mock requests
type MockHandler struct {
	e       *echo.Echo
	service *services.MockService
	log     *logrus.Logger
}

// NewMockHandler creates new MockHandler
func NewMockHandler(e *echo.Echo, service *services.MockService, log *logrus.Logger) *MockHandler {
	return &MockHandler{e: e, service: service, log: log}
}

// RegisterRoutes registers routes for MockHandler
func (h *MockHandler) RegisterRoutes() {
	h.e.Any("/mock/*", h.any)
}

func (h *MockHandler) any(c echo.Context) error {
	request := c.Request()
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to read request body", Details: err.Error()})
	}

	path := "/" + c.Param("*")
	mock, err := h.service.MatchMock(request.Method, path, string(bodyBytes), request.Header, request.Cookies(), request.URL.Query())
	if err != nil {
		return c.JSON(http.StatusNotFound, MessageNotFound)
	}

	for k, v := range mock.Response.Headers {
		c.Response().Header().Set(k, v)
	}
	for k, v := range mock.Response.Cookies {
		cookie := new(http.Cookie)
		cookie.Name = k
		cookie.Value = v
		c.SetCookie(cookie)
	}
	return c.String(mock.Response.Status, mock.Response.Body)
}
