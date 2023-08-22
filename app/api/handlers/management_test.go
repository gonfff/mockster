package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gonfff/mockster/app/models"
	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMockPingBytes = []byte(`
	{
	"name": "ping",
	"path": "/ping",
	"method": "GET",
	"response": {
		"status": 200
	}
}`)

var testMockPing = &models.Mock{
	Name:   "ping",
	Path:   "/ping",
	Method: "GET",
	Response: models.Response{
		Status: 200,
	},
}

func Test_GetMocks_500(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodGet, "/management/mocks", nil)
	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("GetMocks").Return([]*models.Mock{testMock}, errors.New("test"))
	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func Test_GetMocks_200(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodGet, "/management/mocks", nil)
	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("GetMocks").Return([]*models.Mock{testMock}, nil)
	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_CreateMock_201(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPost, "/management/mocks", bytes.NewBuffer(testMockPingBytes))

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()
	req.Header.Set("Content-Type", "application/json")

	repo.On("AddMock", testMockPing).Return(nil)
	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func Test_CreateMock_BadJson400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPost, "/management/mocks", bytes.NewBuffer([]byte("test")))

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()
	req.Header.Set("Content-Type", "application/json")

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_CreateMock_InvalidMock400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPost, "/management/mocks", bytes.NewBuffer([]byte("{}")))

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()
	req.Header.Set("Content-Type", "application/json")

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_CreateMock_RepoFailed400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPost, "/management/mocks", bytes.NewBuffer(testMockPingBytes))

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()
	req.Header.Set("Content-Type", "application/json")

	repo.On("AddMock", mock.Anything).Return(errors.New("test"))
	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_DeleteMock_RepoFailed400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodDelete, "/management/mocks/test", nil)

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("DeleteMock", mock.Anything).Return(errors.New("test"))
	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_DeleteMock_OK(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodDelete, "/management/mocks/test", nil)

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("DeleteMock", mock.Anything).Return(nil)
	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_UpdateMock_BadJson400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPut, "/management/mocks/test", bytes.NewBuffer([]byte("test")))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_UpdateMock_Invalid400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPut, "/management/mocks/test", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_UpdateMock_DeleteFailed400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPut, "/management/mocks/test", bytes.NewBuffer(testMockPingBytes))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("DeleteMock", mock.Anything).Return(errors.New("test"))

	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_UpdateMock_AddFailed400(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPut, "/management/mocks/test", bytes.NewBuffer(testMockPingBytes))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("DeleteMock", mock.Anything).Return(nil)
	repo.On("AddMock", mock.Anything).Return(errors.New("test"))

	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_UpdateMock_OK(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodPut, "/management/mocks/test", bytes.NewBuffer(testMockPingBytes))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("DeleteMock", mock.Anything).Return(nil)
	repo.On("AddMock", mock.Anything).Return(nil)

	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_ExportMocks_RepoErr500(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodGet, "/management/mocks/actions/export", bytes.NewBuffer(testMockPingBytes))

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("GetMocks", mock.Anything).Return([]*models.Mock{}, errors.New("test"))

	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func Test_ExportMocks_OK(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}

	req := httptest.NewRequest(http.MethodGet, "/management/mocks/actions/export", bytes.NewBuffer(testMockPingBytes))

	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	h := NewManagementHandler(e, repo, log)
	h.RegisterRoutes()

	repo.On("GetMocks", mock.Anything).Return([]*models.Mock{}, nil)

	e.ServeHTTP(rec, req)
	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, rec.Code)
}
