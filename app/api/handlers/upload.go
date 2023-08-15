package handlers

import (
	"io"
	"net/http"

	"github.com/gonfff/mockster/app/parsers"
	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
)

// UploadHandler handles upload file requests
type UploadHandler struct {
	e    *echo.Echo
	repo repository.MockRepository
}

// NewUploadHandler creates new UploadHandler
func NewUploadHandler(e *echo.Echo, repo repository.MockRepository) *UploadHandler {
	return &UploadHandler{e: e, repo: repo}
}

// RegisterRoutes registers routes for UploadHandler
func (h *UploadHandler) RegisterRoutes() {
	h.e.POST("/load-yaml", h.LoadYAML)
}

// LoadYAML loads mocks from yaml file
func (h *UploadHandler) LoadYAML(c echo.Context) error {

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to read request body"})
	}

	mocks, err := parsers.ParseYAML(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to parse mocks"})
	}

	errs := make([]string, 0)

	err = h.repo.DeleteAllMocks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Message{Message: "Failed to delete mocks", Details: []string{err.Error()}})
	}
	for _, mock := range mocks {
		// because address of range variable is reused
		newMock := mock
		err = h.repo.AddMock(&newMock)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return c.JSON(http.StatusMultiStatus, Message{Message: "Partial failed to add mocks", Details: errs})
	}
	return c.JSON(http.StatusCreated, MessageSuccess)
}
