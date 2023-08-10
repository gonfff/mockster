package app

import (
	"errors"
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/sirupsen/logrus"
)

// Config is the application configuration. It is filled by the environment
type Config struct {
	Environment      string `env:"ENVIRONMENT,notEmpty"`
	DisableGreetings bool   `env:"DISABLE_GREETINGS" envDefault:"true"`
	InitFilePath     string `env:"INIT_FILE_PATH"`
	LogFormatter     string `env:"LOG_FORMATTER" envDefault:"text"`
	LogLevel         string `env:"LOG_LEVEL" envDefault:"info"`
	RecreateDB       bool   `env:"RECREATE_DB"`
	Port             int    `env:"PORT" envDefault:"8080"`

	logLevel logrus.Level
}

// validateEnvironment validates the ENVIRONMENT variable
// possible values: local, development, production
func (c *Config) validateEnvironment() error {
	envs := map[string]string{
		"local":       "",
		"developnemt": "",
		"production":  "",
	}

	if _, ok := envs[c.Environment]; !ok {
		return errors.New("value must be one of: local, development, production")
	}
	return nil
}

// validateLogFormatter validates the LOG_FORMATTER variable
// possible values: text, json
func (c *Config) validateLogFormatter() error {
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
func (c *Config) validateLogLevel() error {
	lvl, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}
	c.logLevel = lvl
	return nil
}

// NewConfig creates a new Config instance and fills it with the environment variables
// It also validates the environment variables
// If any validation fails, the application will exit
func NewConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	if err := cfg.validateEnvironment(); err != nil {
		log.Fatal(err)
	}
	if err := cfg.validateLogFormatter(); err != nil {
		log.Fatal(err)
	}
	if err := cfg.validateLogLevel(); err != nil {
		log.Fatal(err)
	}

	return cfg
}
