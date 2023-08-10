package app

import (
	"os"

	"github.com/gonfff/mockster/app/api/handlers"
	"github.com/gonfff/mockster/app/api/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func initLogger() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func registerRoutes(app *echo.Echo) {
	app.GET("/ping", handlers.Ping)

}

func registerMiddlewares(app *echo.Echo) {
	app.Use(middlewares.AccessLogMiddleware(logger))
	app.Use(middlewares.RecoverMiddleware(logger))
	app.Pre(middleware.RemoveTrailingSlash())
}

func RunApp() {
	app := echo.New()

	app.HideBanner = true
	app.HidePort = true

	initLogger()
	registerRoutes(app)
	registerMiddlewares(app)
	logger.Info("Application started")
	logger.Info("Listening on port 8080")
	app.Start(":8080")
}
