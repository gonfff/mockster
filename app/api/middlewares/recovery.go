package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// RecoverMiddleware is the middleware that recovers from panics and logs the error
// It returns a middleware function that can be used to register the middleware
// It uses the provided logger to log the error
var RecoverMiddleware = func(log *logrus.Logger) echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(
		middleware.RecoverConfig{
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				log.WithFields(logrus.Fields{
					"error":  err.Error(),
					"stack":  string(stack),
					"URI":    c.Request().RequestURI,
					"method": c.Request().Method,
					"status": c.Response().Status,
					"type":   "recovery",
				}).Error("recovery")
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			},
		})
}
