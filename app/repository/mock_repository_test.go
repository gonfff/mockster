package repository

import (
	"testing"

	"github.com/gonfff/mockster/app/configs"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitRepository_InMemory(t *testing.T) {
	logger := logrus.New()
	cfg := &configs.AppConfig{
		StorageType: "in_memory",
	}

	repo, err := InitRepository(cfg, logger)
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}

func TestInitRepository_Default(t *testing.T) {
	logger := logrus.New()
	cfg := &configs.AppConfig{
		StorageType: "unknown_type",
	}

	repo, err := InitRepository(cfg, logger)
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}
