package repository

import (
	"github.com/gonfff/mockster/app/models"
	"github.com/stretchr/testify/mock"
)

type TestRepository struct {
	mock.Mock
}

func (m *TestRepository) GetMock(name string) (*models.Mock, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Mock), args.Error(1)
}

func (m *TestRepository) GetMocks() ([]*models.Mock, error) {
	args := m.Called()
	return args.Get(0).([]*models.Mock), args.Error(1)
}

func (m *TestRepository) AddMock(mock *models.Mock) error {
	args := m.Called(mock)
	return args.Error(0)
}

func (m *TestRepository) DeleteMock(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *TestRepository) DeleteAllMocks() error {
	args := m.Called()
	return args.Error(0)
}

func (m *TestRepository) ChangeName(oldName string, newName string) error {
	args := m.Called(oldName, newName)
	return args.Error(0)
}

func (m *TestRepository) GetMockNames(endpoint string) ([]string, error) {
	args := m.Called(endpoint)
	return args.Get(0).([]string), args.Error(1)
}
