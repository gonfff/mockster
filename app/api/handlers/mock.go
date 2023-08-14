package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
)

// MockHandler is the handler for the /mock/:path route
type MockHandler struct {
	repo repository.MockRepository
}

// NewMockHandler creates a new MockHandler
func NewMockHandler(repo repository.MockRepository) *MockHandler {
	return &MockHandler{repo: repo}
}

// RegisterRoutes registers the routes for the handler
func (h *MockHandler) RegisterRoutes(e *echo.Echo) {
	e.Any("/mock/:path", h.Mock)
}

// Mock is the handler for the /mock/:path route
func (h *MockHandler) Mock(c echo.Context) error {
	// prepare
	path := "/" + c.Param("path")
	request := c.Request()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}

	bodyStr := string(body)
	bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
	bodyStr = strings.ReplaceAll(bodyStr, "\t", "")

	name, err := h.repo.GetNameByPathMethod(request.Method, path, bodyStr)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not found or body not provided or wrong body"})
	}
	// get mock
	mock, err := h.repo.GetMock(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error() + " " + name})
	}

	// validate request against mock
	for k, v := range mock.Request.Headers {
		if request.Header.Get(k) != v {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Header " + k + " is not equal to " + v})
		}
	}

	for k, v := range mock.Request.Cookies {
		cookie, err := request.Cookie(k)
		if err != nil || cookie.Value != v {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Cookie " + k + " is not equal to " + v})
		}
	}

	for k, v := range mock.Request.QueryParams {
		if request.URL.Query().Get(k) != v {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "QueryParam " + k + " is not equal to " + v})
		}
	}

	if mock.Request.Body != bodyStr {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Body is not equal to " + mock.Request.Body})
	}

	// build response from mock
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
