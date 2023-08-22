package cmd

import (
	"os"
	"testing"

	"github.com/gonfff/mockster/app/configs"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	cfg := &configs.AppConfig{
		LogFormatter: "json",
		IntLogLevel:  logrus.InfoLevel,
	}

	logger := newLogger(cfg)

	assert.NotNil(t, logger)
	assert.IsType(t, &logrus.Logger{}, logger)
	assert.IsType(t, &logrus.JSONFormatter{}, logger.Formatter)
	assert.Equal(t, os.Stdout, logger.Out)
	assert.Equal(t, logrus.InfoLevel, logger.Level)
}
