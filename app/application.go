package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gonfff/mockster/app/api/handlers"
	"github.com/gonfff/mockster/app/api/middlewares"
	"github.com/gonfff/mockster/app/configs"
	"github.com/gonfff/mockster/app/parsers"
	"github.com/gonfff/mockster/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// App is the main application
type App struct {
	e    *echo.Echo
	cfg  *configs.AppConfig
	log  *logrus.Logger
	repo repository.MockRepository
}

// NewApp creates a new application
func NewApp() *App {
	// get configuration
	cfg, err := configs.NewConfig()
	if err != nil {
		log := logrus.New()
		log.WithError(err).Fatal("Failed while reading configuration")
	}

	// initialize logger
	log := newLogger(cfg)

	// initialize repository
	repo, err := repository.InitRepository(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize repository")
	}

	e := echo.New()
	return &App{
		e:    e,
		cfg:  cfg,
		log:  log,
		repo: repo,
	}
}

// newLogger initializes the application logger
func newLogger(cfg *configs.AppConfig) *logrus.Logger {
	log := logrus.New()

	if cfg.LogFormatter == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(cfg.IntLogLevel)
	return log
}

// Setup sets up the application
func (app *App) Setup() {
	app.e.HideBanner = app.cfg.DisableGreetings
	app.e.HidePort = app.cfg.DisableGreetings

	handlers.NewPingHandler(app.e).RegisterRoutes()
	handlers.NewMockHandler(app.e, app.repo, app.log).RegisterRoutes()
	handlers.NewUploadHandler(app.e, app.repo).RegisterRoutes()

	app.registerMiddlewares()
	app.loadInitialMocks()
}

// registerMiddlewares registers the application middlewares
func (app *App) registerMiddlewares() {
	app.e.Use(middlewares.AccessLogMiddleware(app.log))
	app.e.Use(middlewares.RecoverMiddleware(app.log))
	app.e.Pre(middleware.RemoveTrailingSlash())
}

// loadInitialMocks loads the initial mocks from the mock file
func (app *App) loadInitialMocks() {
	if app.cfg.MockFilePath == "" {
		return
	}
	mocks, err := parsers.FileYAML(app.cfg.MockFilePath)
	if err != nil {
		app.log.WithError(err).Error("Failed to parse mocks")
		return
	}

	for _, mock := range mocks {
		// because address of range variable is reused
		newMock := mock
		err = app.repo.AddMock(&newMock)
		if err != nil {
			app.log.WithError(err).Error("Failed to add mock")
		}
	}
}

// Start starts the application
func (app *App) Start() {
	app.log.Info("Application started")
	app.log.Info("Listening on port 8080")
	err := app.e.Start(fmt.Sprintf(":%v", app.cfg.Port))
	if err != nil && err != http.ErrServerClosed {
		app.log.WithError(err).Fatal("Application failed")
	}
}
