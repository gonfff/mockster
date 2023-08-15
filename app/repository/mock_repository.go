package repository

import (
	"github.com/gonfff/mockster/app/configs"
	"github.com/gonfff/mockster/app/models"
	"github.com/sirupsen/logrus"
)

// MockRepository is the interface for the mock repository
type MockRepository interface {
	GetMock(name string) (*models.Mock, error)
	GetMocks() ([]*models.Mock, error)
	AddMock(mock *models.Mock) error
	UpdateMock(mock *models.Mock) error
	DeleteMock(name string) error
	DeleteAllMocks() error
	ChangeName(oldName string, newName string) error
	GetMockNames(endpoint string) ([]string, error)
	// todo add methods for exporting and importing mocks from yaml
}

// InitRepository initializes the repository based on the configuration
func InitRepository(cfg *configs.AppConfig, log *logrus.Logger) (MockRepository, error) {
	var repo MockRepository
	switch cfg.StorageType {
	case "in_memory":
		repo = NewInMemoryRepository(log)
	default:
		repo = NewInMemoryRepository(log)
	}
	return repo, nil
}
