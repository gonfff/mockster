package repository

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gonfff/mockster/app/models"
	"github.com/sirupsen/logrus"
)

// NewInMemoryRepository creates a new InMemoryRepository
func NewInMemoryRepository(log *logrus.Logger) *InMemoryRepository {
	r := &InMemoryRepository{log: log}
	if r.storage == nil {
		r.storage = make(map[string]*models.Mock)
	}
	if r.pathIndex == nil {
		r.pathIndex = make(map[string]string)
	}
	return r

}

// InMemoryRepository is an in-memory implementation of the MockRepository
type InMemoryRepository struct {
	log     *logrus.Logger
	mu      sync.RWMutex
	storage map[string]*models.Mock

	order     []string
	pathIndex map[string]string
}

// GetNameByPathMethod returns the name of the mock with the given path and method
func (r *InMemoryRepository) GetNameByPathMethod(method, path, body string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// complex composite key for greater uniqueness on partial matches
	key := fmt.Sprintf("%v %v %v", method, path, strings.ReplaceAll(body, "\n", ""))
	name, ok := r.pathIndex[key]
	if !ok {
		return "", errors.New("mock with this path and method does not exist")
	}
	return name, nil
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
		return errors.New("mock with this name already exists")
	}

	r.storage[mock.Name] = mock
	key := fmt.Sprintf("%v %v %v", mock.Method, mock.Path, mock.Request.Body)
	r.pathIndex[key] = mock.Name
	r.order = append(r.order, mock.Name)
	sort.Strings(r.order)

	return nil
}

// UpdateMock updates the mock with the given name
func (r *InMemoryRepository) UpdateMock(mock *models.Mock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[mock.Name]; !ok {
		return errors.New("mock with this name does not exist")
	}

	r.storage[mock.Name] = mock
	key := fmt.Sprintf("%v %v %v", mock.Method, mock.Path, mock.Response.Body)
	r.pathIndex[key] = mock.Name
	return nil
}

// DeleteMock deletes the mock with the given name
func (r *InMemoryRepository) DeleteMock(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[name]; !ok {
		return errors.New("mock with this name does not exist")
	}
	delete(r.storage, name)
	return nil
}

// ChangeName changes the name of the mock
func (r *InMemoryRepository) ChangeName(oldName, newName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	mock, ok := r.storage[oldName]
	if !ok {
		return errors.New("mock with this name does not exist")
	}

	r.storage[newName] = r.storage[oldName]
	key := fmt.Sprintf("%v %v %v", mock.Method, mock.Path, mock.Response.Body)
	r.pathIndex[key] = newName
	delete(r.storage, oldName)
	return nil
}
