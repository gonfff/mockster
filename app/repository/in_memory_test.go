package repository

import (
	"testing"

	"github.com/gonfff/mockster/app/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func mockNewInMemoryRepository() *InMemoryRepository {
	log := logrus.New()
	r := NewInMemoryRepository(log)
	return r
}

func TestGetMock(t *testing.T) {
	r := mockNewInMemoryRepository()
	mock := &models.Mock{Name: "mock1"}
	r.storage[mock.Name] = mock

	testCases := []struct {
		name     string
		mockName string
		err      bool
	}{
		{"existing mock", "mock1", false},
		{"non-existing mock", "mock2", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			received, err := r.GetMock(tc.mockName)
			if tc.err {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if mock != received {
					t.Errorf("expected %v, got %v", mock, mock)
				}
			}
		})
	}
}

func TestGetMocks(t *testing.T) {
	r := mockNewInMemoryRepository()
	mock1 := &models.Mock{Name: "mock1"}
	mock2 := &models.Mock{Name: "mock2"}
	r.storage[mock1.Name] = mock1
	r.storage[mock2.Name] = mock2
	r.order = append(r.order, mock1.Name, mock2.Name)

	mocks, err := r.GetMocks()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(mocks))
	assert.Equal(t, mock1, mocks[0])
	assert.Equal(t, mock2, mocks[1])
}

func TestAddMock(t *testing.T) {
	r := mockNewInMemoryRepository()
	newMock := &models.Mock{Name: "newMock"}
	existingMock := &models.Mock{Name: "existingMock"}
	r.storage[existingMock.Name] = newMock
	r.order = append(r.order, existingMock.Name)

	testCases := []struct {
		name       string
		mock       *models.Mock
		err        bool
		storageLen int
	}{
		{"new mock", newMock, false, 2},
		{"existing mock", existingMock, true, 2},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				err := r.AddMock(tc.mock)
				if tc.err {
					assert.Error(t, err)
				}
				assert.Equal(t, tc.storageLen, len(r.storage))
			})
	}
}

func TestDeleteMock(t *testing.T) {
	r := mockNewInMemoryRepository()
	notExistingMock := &models.Mock{Name: "notExistingMock"}
	existingMock := &models.Mock{Name: "existingMock"}
	r.storage[existingMock.Name] = existingMock
	r.order = append(r.order, existingMock.Name)

	testCases := []struct {
		name       string
		mockName   string
		err        bool
		storageLen int
	}{
		{"not existing mock", notExistingMock.Name, true, 1},
		{"existing mock", existingMock.Name, false, 0},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				err := r.DeleteMock(tc.mockName)
				if tc.err {
					assert.Error(t, err)
				}
				assert.Equal(t, tc.storageLen, len(r.storage))
			})
	}
}

func TestChangeName(t *testing.T) {
	r := mockNewInMemoryRepository()
	notExistingMock := &models.Mock{Name: "notExistingMock"}
	existingMock := &models.Mock{Name: "existingMock"}
	r.storage[existingMock.Name] = existingMock
	r.order = append(r.order, existingMock.Name)

	testCases := []struct {
		name       string
		oldName    string
		newName    string
		err        bool
		storageLen int
	}{
		{"not existing mock", notExistingMock.Name, "newMock", true, 1},
		{"existing mock", existingMock.Name, "newMock", false, 1},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				err := r.ChangeName(tc.oldName, tc.newName)
				if tc.err {
					assert.Error(t, err)
				}
				assert.Equal(t, tc.storageLen, len(r.storage))
			})
	}
}

func TestDeleteFromEndpoints(t *testing.T) {
	r := mockNewInMemoryRepository()
	mock1 := &models.Mock{Name: "mock1", Method: "GET", Path: "/path1"}
	mock2 := &models.Mock{Name: "mock2", Method: "GET", Path: "/path1"}
	mock3 := &models.Mock{Name: "mock3", Method: "POST", Path: "/path1"}
	r.storage[mock1.Name] = mock1
	r.storage[mock2.Name] = mock2
	r.storage[mock3.Name] = mock3
	r.order = append(r.order, mock1.Name, mock2.Name, mock3.Name)

	r.endpointMocks["GET /path1"] = []string{mock1.Name, mock2.Name}
	r.endpointMocks["POST /path1"] = []string{mock3.Name}

	r.deleteFromEndpoints(mock1)
	assert.Equal(t, 1, len(r.endpointMocks["GET /path1"]))
	assert.Equal(t, 1, len(r.endpointMocks["POST /path1"]))
	assert.Equal(t, 3, len(r.storage))
	assert.Equal(t, 3, len(r.order))

}

func TestGetMockNamesByEndpoint(t *testing.T) {
	r := mockNewInMemoryRepository()
	mock1 := &models.Mock{Name: "mock1", Method: "GET", Path: "/path1"}
	mock2 := &models.Mock{Name: "mock2", Method: "GET", Path: "/path1"}
	mock3 := &models.Mock{Name: "mock3", Method: "POST", Path: "/path1"}
	r.storage[mock1.Name] = mock1
	r.storage[mock2.Name] = mock2
	r.storage[mock3.Name] = mock3
	r.order = append(r.order, mock1.Name, mock2.Name, mock3.Name)
	r.endpointMocks["GET /path1"] = []string{mock1.Name, mock2.Name}
	r.endpointMocks["POST /path1"] = []string{mock3.Name}

	testCases := []struct {
		name     string
		endpoint string
		expected []string
		err      bool
	}{
		{"existing endpoint", "GET /path1", []string{"mock1", "mock2"}, false},
		{"non existing endpoint", "POST /path2", []string{}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			received, err := r.GetMockNames(tc.endpoint)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.expected, received)
			}
		})
	}

}

func TestDeleteAllMocks(t *testing.T) {
	r := mockNewInMemoryRepository()
	mock1 := &models.Mock{Name: "mock1", Method: "GET", Path: "/path1"}
	mock2 := &models.Mock{Name: "mock2", Method: "GET", Path: "/path1"}
	mock3 := &models.Mock{Name: "mock3", Method: "POST", Path: "/path1"}
	r.storage[mock1.Name] = mock1
	r.storage[mock2.Name] = mock2
	r.storage[mock3.Name] = mock3
	r.order = append(r.order, mock1.Name, mock2.Name, mock3.Name)
	r.endpointMocks["GET /path1"] = []string{mock1.Name, mock2.Name}
	r.endpointMocks["POST /path1"] = []string{mock3.Name}

	r.DeleteAllMocks()
	assert.Equal(t, 0, len(r.storage))
	assert.Equal(t, 0, len(r.order))
	assert.Equal(t, 0, len(r.endpointMocks))
}
