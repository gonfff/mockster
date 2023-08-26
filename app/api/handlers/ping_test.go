package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler_RegisterRoutes(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
	rec := httptest.NewRecorder()
	_ = e.NewContext(req, rec)

	handler := NewPingHandler(e)

	handler.RegisterRoutes()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message": "pong"}`, rec.Body.String())
}
