package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRecoverMiddleware(t *testing.T) {
	logger := logrus.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := RecoverMiddleware(logger)

	handler := func(c echo.Context) error {
		panic("test panic")
	}

	_ = middleware(handler)(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "Internal Server Error", rec.Body.String())
}
