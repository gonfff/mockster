package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gonfff/mockster/app/models"
	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type MockHandler struct {
	mu   sync.Mutex
	e    *echo.Echo
	repo repository.MockRepository
	log  *logrus.Logger

	endpointMocks map[string][]string
}

func NewMockHandler(e *echo.Echo, repo repository.MockRepository, log *logrus.Logger) *MockHandler {

	return &MockHandler{e: e, repo: repo, log: log, endpointMocks: make(map[string][]string)}
}

func (h *MockHandler) RegisterRoutes() {
	h.e.Group("/mocks")

	// todo
}

func (h *MockHandler) RegisterMock(mock *models.Mock) {
	h.mu.Lock()
	defer h.mu.Unlock()

	mockNames := h.endpointMocks[mock.Method+mock.Path]
	mockNames = append(mockNames, mock.Name)
	h.endpointMocks[mock.Method+mock.Path] = mockNames

	if len(mockNames) == 1 {
		switch mock.Method {
		case "GET":
			h.e.GET(mock.Path, h.Any)
		case "POST":
			h.e.POST(mock.Path, h.Any)
		case "PUT":
			h.e.PUT(mock.Path, h.Any)
		case "DELETE":
			h.e.DELETE(mock.Path, h.Any)
		case "PATCH":
			h.e.PATCH(mock.Path, h.Any)
		case "HEAD":
			h.e.HEAD(mock.Path, h.Any)
		case "OPTIONS":
			h.e.OPTIONS(mock.Path, h.Any)
		}
	}
}

func (h *MockHandler) UnregisterMock(mock *models.Mock) {
	h.mu.Lock()
	defer h.mu.Unlock()

	mockNames := h.endpointMocks[mock.Method+mock.Path]

	for i, v := range mockNames {
		if v == mock.Name {
			mockNames = append(mockNames[:i], mockNames[i+1:]...)
			break
		}
	}
	h.endpointMocks[mock.Method+mock.Path] = mockNames

	if len(mockNames) == 0 {
		notFoundFunc := func(c echo.Context) error {
			return c.JSON(404, JSONNotFound)
		}
		switch mock.Method {
		case "GET":
			h.e.GET(mock.Path, notFoundFunc)
		case "POST":
			h.e.POST(mock.Path, notFoundFunc)
		case "PUT":
			h.e.PUT(mock.Path, notFoundFunc)
		case "DELETE":
			h.e.DELETE(mock.Path, notFoundFunc)
		case "PATCH":
			h.e.PATCH(mock.Path, notFoundFunc)
		case "HEAD":
			h.e.HEAD(mock.Path, notFoundFunc)
		case "OPTIONS":
			h.e.OPTIONS(mock.Path, notFoundFunc)
		}
	}
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
			return errors.New(fmt.Sprintf("Header %v is not equal to %v", k, v))
		}
	}

	for k, v := range mock.Request.Cookies {
		cookie, err := request.Cookie(k)
		if err != nil || cookie.Value != v {
			return errors.New(fmt.Sprintf("Cookie %v is not equal to %v", k, v))
		}
	}

	for k, v := range mock.Request.QueryParams {
		if request.URL.Query().Get(k) != v {
			return errors.New(fmt.Sprintf("QueryParam %v is not equal to %v", k, v))
		}
	}

	if mock.Request.Body != bodyStr {
		return errors.New(fmt.Sprintf("Body is not equal to %v", mock.Request.Body))
	}

	return nil
}

func (h *MockHandler) Any(c echo.Context) error {
	mockNames := h.endpointMocks[mock.Method+mock.Path]

	if len(mockNames) == 0 {
		return c.JSON(404, JSONNotFound)
	}

	for i, mockName := range mockNames {
		mock, err := h.repo.GetMock(mockName)
		if err != nil && i == len(mockNames)-1 {
			return c.JSON(400, JSONMessageError(err))
		}
		c.validateQuery(c, mock)
		if err != nil {
			return c.JSON(400, JSONMessageError(err))
		}

		//build response
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
}
