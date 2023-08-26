package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAccessLogMiddleware(t *testing.T) {
	logger := logrus.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := AccessLogMiddleware(logger)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
	err := middleware(handler)(c)

	assert.Nil(t, err)

}
