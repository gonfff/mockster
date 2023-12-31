package configs

import (
	"errors"

	"github.com/caarlos0/env/v9"
	"github.com/sirupsen/logrus"
)

// AppConfig is the application configuration. It is filled by the environment
type AppConfig struct {
	DisableGreetings bool   `env:"DISABLE_GREETINGS" envDefault:"true"`
	MockFilePath     string `env:"MOCK_FILE_PATH"`
	LogFormatter     string `env:"LOG_FORMATTER" envDefault:"text"`
	LogLevel         string `env:"LOG_LEVEL" envDefault:"info"`
	RecreateDB       bool   `env:"RECREATE_DB"`
	Port             int    `env:"PORT" envDefault:"8080"`
	StorageType      string `env:"STORAGE" envDefault:"in_memory"`
	StaticPath       string `env:"STATIC_PATH" envDefault:"app/static"`

	IntLogLevel logrus.Level
}

// validateLogFormatter validates the LOG_FORMATTER variable
// possible values: text, json
func (c *AppConfig) validateLogFormatter() error {
	logFormatters := map[string]string{
		"text": "",
		"json": "",
	}

	if _, ok := logFormatters[c.LogFormatter]; !ok {
		return errors.New("value must be one of: text, json")
	}
	return nil
}

// validateLogLevel validates the LOG_LEVEL variable
// possible values: panic, fatal, error, warn, info, debug, trace
func (c *AppConfig) validateLogLevel() error {
	lvl, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}
	c.IntLogLevel = lvl
	return nil
}

// NewConfig creates a new AppConfig instance and fills it with the environment variables
// It also validates the environment variables
// If any validation fails, the application will exit
func NewConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if err := cfg.validateLogFormatter(); err != nil {
		return nil, err
	}
	if err := cfg.validateLogLevel(); err != nil {
		return nil, err
	}

	return cfg, nil
}
