package repository

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/gonfff/mockster/app/models"
	"github.com/sirupsen/logrus"
)

// InMemoryRepository is an in-memory implementation of the MockRepository
type InMemoryRepository struct {
	log     *logrus.Logger
	mu      sync.RWMutex
	storage map[string]*models.Mock

	order         []string
	endpointMocks map[string][]string
}

// NewInMemoryRepository creates a new InMemoryRepository
func NewInMemoryRepository(log *logrus.Logger) *InMemoryRepository {
	r := &InMemoryRepository{log: log}
	if r.storage == nil {
		r.storage = make(map[string]*models.Mock)
	}
	if r.endpointMocks == nil {
		r.endpointMocks = make(map[string][]string)
	}
	return r

}

// GetMock returns the mock with the given name
func (r *InMemoryRepository) GetMock(name string) (*models.Mock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mock, ok := r.storage[name]
	if !ok {
		return nil, errors.New("mock with this name does not exist")
	}
	return mock, nil
}

// GetMocks returns all mocks
func (r *InMemoryRepository) GetMocks() ([]*models.Mock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mocks := make([]*models.Mock, 0, len(r.storage))
	for _, name := range r.order {
		mocks = append(mocks, r.storage[name])
	}
	return mocks, nil
}

// AddMock adds a new mock
func (r *InMemoryRepository) AddMock(mock *models.Mock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[mock.Name]; ok {
		return fmt.Errorf("mock with name \"%s\" already exists", mock.Name)
	}

	r.storage[mock.Name] = mock
	key := fmt.Sprintf("%v %v", mock.Method, mock.Path)
	r.endpointMocks[key] = append(r.endpointMocks[key], mock.Name)
	r.order = append(r.order, mock.Name)
	sort.Strings(r.order)

	return nil
}

// DeleteMock deletes the mock with the given name
func (r *InMemoryRepository) DeleteMock(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	mock, ok := r.storage[name]
	if !ok {
		return fmt.Errorf("mock with name \"%s\" does not exist", name)
	}

	r.deleteFromEndpoints(mock)
	delete(r.storage, name)
	return nil
}

// ChangeName changes the name of the mock
func (r *InMemoryRepository) ChangeName(oldName, newName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	mock, ok := r.storage[oldName]
	if !ok {
		return fmt.Errorf("mock with name \"%s\" does not exist", oldName)

	}
	r.deleteFromEndpoints(mock)

	mock.Name = newName
	endpoint := fmt.Sprintf("%v %v", mock.Method, mock.Path)
	r.endpointMocks[endpoint] = append(r.endpointMocks[endpoint], newName)
	delete(r.storage, oldName)
	r.storage[newName] = mock
	return nil
}

// deleteFromEndpoints deletes the mock from the endpointMocks map
func (r *InMemoryRepository) deleteFromEndpoints(mock *models.Mock) {
	endpoint := fmt.Sprintf("%v %v", mock.Method, mock.Path)
	for i, name := range r.endpointMocks[endpoint] {
		if name == mock.Name {
			r.endpointMocks[endpoint] = append(r.endpointMocks[endpoint][:i], r.endpointMocks[endpoint][i+1:]...)
			break
		}
	}
}

// GetMockNames returns all mock names for the given endpoint (method + path)
func (r *InMemoryRepository) GetMockNames(endpoint string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mockNames, ok := r.endpointMocks[endpoint]
	if !ok {
		return nil, errors.New("endpoint does not exist")
	}
	return mockNames, nil
}

// DeleteAllMocks deletes all mocks
func (r *InMemoryRepository) DeleteAllMocks() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage = make(map[string]*models.Mock)
	r.order = make([]string, 0)
	r.endpointMocks = make(map[string][]string)
	return nil
}
