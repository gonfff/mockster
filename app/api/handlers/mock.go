package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gonfff/mockster/app/models"
	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// MockHandler handles mock requests
type MockHandler struct {
	e    *echo.Echo
	repo repository.MockRepository
	log  *logrus.Logger
}

// NewMockHandler creates new MockHandler
func NewMockHandler(e *echo.Echo, repo repository.MockRepository, log *logrus.Logger) *MockHandler {
	return &MockHandler{e: e, repo: repo, log: log}
}

// RegisterRoutes registers routes for MockHandler
func (h *MockHandler) RegisterRoutes() {
	h.e.Any("/mock/:path", h.any)
}

func (h *MockHandler) validateQuery(c echo.Context, mock *models.Mock) error {
	request := c.Request()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	bodyStr := string(body)
	bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
	bodyStr = strings.ReplaceAll(bodyStr, "\t", "")

	for k, v := range mock.Request.Headers {
		if request.Header.Get(k) != v {
			return fmt.Errorf("header %v is not equal to %v", k, v)
		}
	}

	for k, v := range mock.Request.Cookies {
		cookie, err := request.Cookie(k)
		if err != nil || cookie.Value != v {
			return fmt.Errorf("cookie %v is not equal to %v", k, v)
		}
	}

	for k, v := range mock.Request.QueryParams {
		if request.URL.Query().Get(k) != v {
			return fmt.Errorf("QueryParam %v is not equal to %v", k, v)
		}
	}

	if mock.Request.Body != bodyStr {
		return fmt.Errorf("body is not equal to %v", mock.Request.Body)
	}

	return nil
}

func (h *MockHandler) any(c echo.Context) error {
	// prepare
	request := c.Request()
	endpoint := fmt.Sprintf("%v /%v", request.Method, c.Param("path"))
	mockNames, err := h.repo.GetMockNames(endpoint)

	// check if mock exists
	if err != nil {
		return c.JSON(404, JSONMessageError(err))
	}
	if len(mockNames) == 0 {
		return c.JSON(500, JSONMessageText("No mocks found for this endpoint"))
	}

	// validate
	for i, mockName := range mockNames {
		mock, err := h.repo.GetMock(mockName)
		if err != nil && i == len(mockNames)-1 {
			return c.JSON(400, JSONMessageError(err))
		} else if err != nil {
			continue
		}

		err = h.validateQuery(c, mock)

		if err != nil && i == len(mockNames)-1 {
			return c.JSON(400, JSONMessageError(err))
		} else if err != nil {
			continue
		}

		// build response
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
	return c.JSON(500, JSONMessageText("Somthing went wrong"))
}
