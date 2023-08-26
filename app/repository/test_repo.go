package repository

import (
	"github.com/gonfff/mockster/app/models"
	libmock "github.com/stretchr/testify/mock"
)

// TestRepository is a mock repository for testing
type TestRepository struct {
	libmock.Mock
}

// GetMock returns a mock by name
func (m *TestRepository) GetMock(name string) (*models.Mock, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Mock), args.Error(1)
}

// GetMocks returns all mocks
func (m *TestRepository) GetMocks() ([]*models.Mock, error) {
	args := m.Called()
	return args.Get(0).([]*models.Mock), args.Error(1)
}

// AddMock adds a new mock
func (m *TestRepository) AddMock(mock *models.Mock) error {
	args := m.Called(mock)
	return args.Error(0)
}

// DeleteMock deletes a mock
func (m *TestRepository) DeleteMock(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

// DeleteAllMocks deletes all mocks
func (m *TestRepository) DeleteAllMocks() error {
	args := m.Called()
	return args.Error(0)
}

// ChangeName changes the name of a mock
func (m *TestRepository) ChangeName(oldName, newName string) error {
	args := m.Called(oldName, newName)
	return args.Error(0)
}

// GetMockNames returns all mock names
func (m *TestRepository) GetMockNames(endpoint string) ([]string, error) {
	args := m.Called(endpoint)
	return args.Get(0).([]string), args.Error(1)
}
