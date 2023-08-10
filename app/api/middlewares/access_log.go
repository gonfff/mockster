package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

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
