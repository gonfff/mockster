/*
Package app is the main application package.
*/
package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gonfff/mockster/app/api/handlers"
	"github.com/gonfff/mockster/app/api/middlewares"
	"github.com/gonfff/mockster/app/parsers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// logger is the application logger
var logger *logrus.Logger

// initLogger initializes the application logger
// configuration is read from the environment
func initLogger(cfg *Config) {
	logger = logrus.New()

	if cfg.LogFormatter == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(cfg.logLevel)
}

// registerRoutes registers the application routes
func registerRoutes(app *echo.Echo) {
	app.GET("/ping", handlers.Ping)

}

// registerMiddlewares registers the application middlewares
func registerMiddlewares(app *echo.Echo) {
	app.Use(middlewares.AccessLogMiddleware(logger))
	app.Use(middlewares.RecoverMiddleware(logger))
	app.Pre(middleware.RemoveTrailingSlash())
}

// RunApp runs the main application
func RunApp() {
	app := echo.New()
	cfg, err := newConfig()

	initLogger(cfg)

	if err != nil {
		logger.WithError(err).Fatal("Failed while reading configuration")
	}

	mocks, err := parsers.ParseYAML(cfg.MockFilePath)
	if err != nil {
		logger.WithError(err).Error("Failed to parse mocks")
	}
	// todo remove this
	fmt.Println(mocks)

	app.HideBanner = cfg.DisableGreetings
	app.HidePort = cfg.DisableGreetings

	registerRoutes(app)
	registerMiddlewares(app)

	logger.Info("Application started")
	logger.Info("Listening on port 8080")

	err = app.Start(fmt.Sprintf(":%v", cfg.Port))
	if err != nil && err != http.ErrServerClosed {
		logger.WithError(err).Fatal("Application failed")
	}

}
