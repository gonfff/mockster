package repository

import (
	"github.com/gonfff/mockster/app/models"
)

// MockRepository is the interface for the mock repository
type MockRepository interface {
	GetMock(name string) (*models.Mock, error)
	GetMocks() ([]*models.Mock, error)
	AddMock(mock *models.Mock) error
	DeleteMock(name string) error
	UpdateMock(name string, mock *models.Mock) error
	ReplaceAll(mocks []*models.Mock) error
	GetMockNames(endpoint string) ([]string, error)
}

// InitRepository initializes the repository based on the configuration
func InitRepository() (MockRepository, error) {
	return NewInMemoryRepository(), nil
}
