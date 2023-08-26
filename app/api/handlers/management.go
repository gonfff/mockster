package handlers

import (
	"io"
	"net/http"

	"github.com/gonfff/mockster/app/models"
	"github.com/gonfff/mockster/app/parsers"
	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// ManagementHandler is the handler for the management API
type ManagementHandler struct {
	e    *echo.Echo
	repo repository.MockRepository
	log  *logrus.Logger
}

// NewManagementHandler creates a new ManagementHandler
func NewManagementHandler(e *echo.Echo, repo repository.MockRepository, log *logrus.Logger) *ManagementHandler {
	return &ManagementHandler{e: e, repo: repo, log: log}
}

// RegisterRoutes registers the routes for the handler
func (h *ManagementHandler) RegisterRoutes() {
	g := h.e.Group("/management")
	g.GET("/mocks", h.GetMocks)
	g.POST("/mocks", h.CreateMock)
	g.DELETE("/mocks/:name", h.DeleteMock)
	g.PUT("/mocks/:name", h.UpdateMock)
	g.GET("/mocks/actions/export", h.ExportMocks)
	g.POST("/mocks/actions/import", h.ImportMocks)
}

// GetMocks returns all mocks
func (h *ManagementHandler) GetMocks(c echo.Context) error {
	mocks, err := h.repo.GetMocks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Message{Message: "Failed to get mocks", Details: err.Error()})
	}
	return c.JSON(http.StatusOK, PayloadMocks{Items: mocks})

}

// CreateMock creates a new mock
func (h *ManagementHandler) CreateMock(c echo.Context) error {
	mock := new(models.Mock)
	if err := c.Bind(mock); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Invalid JSON data"})
	}

	if err := mock.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Invalid mock", Details: err.Error()})
	}

	if err := h.repo.AddMock(mock); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to add mock", Details: err.Error()})
	}
	return c.JSON(http.StatusCreated, MessageSuccess)
}

// DeleteMock deletes a mock
func (h *ManagementHandler) DeleteMock(c echo.Context) error {
	name := c.Param("name")
	if err := h.repo.DeleteMock(name); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to delete mock", Details: err.Error()})
	}

	return c.JSON(http.StatusOK, MessageSuccess)
}

// UpdateMock updates a mock
func (h *ManagementHandler) UpdateMock(c echo.Context) error {
	name := c.Param("name")
	mock := new(models.Mock)
	if err := c.Bind(mock); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Invalid JSON data"})
	}

	if err := mock.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Invalid mock", Details: err.Error()})
	}

	// todo add atomicity
	if err := h.repo.DeleteMock(name); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to delete mock", Details: err.Error()})
	}

	if err := h.repo.AddMock(mock); err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to add mock", Details: err.Error()})
	}
	return c.JSON(http.StatusOK, MessageSuccess)
}

// ExportMocks exports all mocks to YAML
func (h *ManagementHandler) ExportMocks(c echo.Context) error {
	mocks, err := h.repo.GetMocks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Message{Message: "Failed to get mocks", Details: err.Error()})
	}

	data, err := parsers.ToYAML(mocks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Message{Message: "Failed to convert mocks to YAML", Details: err.Error()})
	}
	c.Response().Header().Set("Content-Disposition", "attachment; filename=mocks.yaml")
	return c.Blob(http.StatusOK, "application/yaml", data)
}

// ImportMocks imports mocks from YAML
func (h *ManagementHandler) ImportMocks(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to read file"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Message{Message: "Failed to open file"})
	}
	defer src.Close()

	data, err := io.ReadAll(src)
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
		m := mock
		err = h.repo.AddMock(&m)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return c.JSON(http.StatusMultiStatus, Message{Message: "Partial failed to add mocks", Details: errs})
	}
	return c.JSON(http.StatusCreated, MessageSuccess)

}
