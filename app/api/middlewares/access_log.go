package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// AccessLogMiddleware is the middleware that logs the access log
// It returns a middleware function that can be used to register the middleware
// It uses the provided logger to log the access log
func AccessLogMiddleware(log *logrus.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency: true,
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"method":  values.Method,
				"URI":     values.URI,
				"status":  values.Status,
				"latency": values.Latency,
				"type":    "access",
			}).Info("access")
			return nil
		},
	})
}
