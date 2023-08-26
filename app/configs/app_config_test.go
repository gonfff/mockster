package configs

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAppConfig_validateLogFormatter(t *testing.T) {
	cases := []struct {
		name         string
		logFormatter string
		expectedErr  bool
	}{
		{
			name:         "ValidLogFormatter",
			logFormatter: "text",
			expectedErr:  false,
		},
		{
			name:         "InvalidLogFormatter",
			logFormatter: "invalid",
			expectedErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := &AppConfig{LogFormatter: tc.logFormatter}
			err := config.validateLogFormatter()
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAppConfig_validateLogLevel(t *testing.T) {
	cases := []struct {
		name        string
		logLevel    string
		expectedErr bool
	}{
		{
			name:        "ValidLogLevel",
			logLevel:    "info",
			expectedErr: false,
		},
		{
			name:        "InvalidLogLevel",
			logLevel:    "invalid",
			expectedErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := &AppConfig{LogLevel: tc.logLevel}
			err := config.validateLogLevel()
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, logrus.InfoLevel, config.IntLogLevel)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name             string
		environmentVars  map[string]string
		expectedErr      bool
		expectedLogLevel string
	}{
		{
			name: "ValidConfig",
			environmentVars: map[string]string{
				"LOG_FORMATTER": "text",
				"LOG_LEVEL":     "info",
			},
			expectedErr:      false,
			expectedLogLevel: "info",
		},
		{
			name:            "InvalidConfig",
			environmentVars: map[string]string{"LOG_LEVEL": "qwe"},
			expectedErr:     true,
		},
		{
			name: "ValidConfig",
			environmentVars: map[string]string{
				"LOG_FORMATTER": "qweqwe",
				"LOG_LEVEL":     "info",
			},
			expectedErr:      true,
			expectedLogLevel: "info",
		},
		{
			name: "ValidConfig",
			environmentVars: map[string]string{
				"LOG_LEVEL": "qweqwe",
			},
			expectedErr:      true,
			expectedLogLevel: "info",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restoreEnv := setEnvVariables(test.environmentVars)
			defer restoreEnv()

			cfg, err := NewConfig()
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedLogLevel, cfg.LogLevel)
			}
		})
	}
}

func setEnvVariables(vars map[string]string) func() {
	originalValues := make(map[string]string)
	for key, value := range vars {
		originalValues[key] = os.Getenv(key)
		os.Setenv(key, value)
	}

	return func() {
		for key, value := range originalValues {
			os.Setenv(key, value)
		}
	}
}
