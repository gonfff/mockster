package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gonfff/mockster/app/models"
	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testMock = &models.Mock{
	Name:   "test",
	Path:   "/test",
	Method: "POST",
	Request: models.Request{
		Headers: map[string]string{
			"Content-Type": "application/text",
		},
		QueryParams: map[string]string{
			"test": "test",
		},
		Cookies: map[string]string{
			"test": "test",
		},
		Body: "test",
	},
	Response: models.Response{
		Status: 200,
		Headers: map[string]string{
			"Content-Type": "application/text",
		},
		Cookies: map[string]string{
			"test": "test",
		},
		Body: "test",
	},
}

func Test_MockHandler_any(t *testing.T) {
	e := echo.New()
	log := logrus.New()
	repo := &repository.TestRepository{}
	req := prepareQuery()
	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	repo.On("GetMock", "test").Return(testMock, nil)
	repo.On("GetMockNames", "POST /test").Return([]string{"test"}, nil)
	h := NewMockHandler(e, repo, log)
	h.RegisterRoutes()

	e.ServeHTTP(rec, req)

	repo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `test`, rec.Body.String())

}

func prepareQuery() *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/mock/test", bytes.NewReader([]byte("test")))
	req.Header.Set("Content-Type", "application/text")
	req.AddCookie(&http.Cookie{Name: "test", Value: "test"})
	q := req.URL.Query()
	q.Add("test", "test")
	req.URL.RawQuery = q.Encode()
	return req
}
