package repository

import (
	"testing"

	"github.com/gonfff/mockster/app/models"
	"github.com/stretchr/testify/assert"
)

func mockNewInMemoryRepository() *InMemoryRepository {
	r := NewInMemoryRepository()
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

func TestUpdateMock(t *testing.T) {
	r := mockNewInMemoryRepository()
	oldMock := &models.Mock{Name: "old", Method: "GET", Path: "/one", Response: models.Response{Status: 200}}
	err := r.AddMock(oldMock)
	assert.NoError(t, err)

	updated := &models.Mock{Name: "new", Method: "POST", Path: "/two", Response: models.Response{Status: 201}}
	err = r.UpdateMock("old", updated)
	assert.NoError(t, err)

	_, err = r.GetMock("old")
	assert.Error(t, err)
	got, err := r.GetMock("new")
	assert.NoError(t, err)
	assert.Equal(t, "POST", got.Method)

	names, err := r.GetMockNames("POST /two")
	assert.NoError(t, err)
	assert.Equal(t, []string{"new"}, names)
}

func TestReplaceAll(t *testing.T) {
	r := mockNewInMemoryRepository()

	mocks := []*models.Mock{
		{Name: "b", Method: "GET", Path: "/b", Response: models.Response{Status: 200}},
		{Name: "a", Method: "POST", Path: "/a", Response: models.Response{Status: 201}},
	}
	err := r.ReplaceAll(mocks)
	assert.NoError(t, err)

	stored, err := r.GetMocks()
	assert.NoError(t, err)
	assert.Len(t, stored, 2)
	assert.Equal(t, "a", stored[0].Name)
	assert.Equal(t, "b", stored[1].Name)

	err = r.ReplaceAll([]*models.Mock{{Name: "x"}, {Name: "x"}})
	assert.Error(t, err)

	stored, err = r.GetMocks()
	assert.NoError(t, err)
	assert.Len(t, stored, 2)
}

func TestUpdateMockAtomicOnConflict(t *testing.T) {
	r := mockNewInMemoryRepository()
	err := r.AddMock(&models.Mock{Name: "one", Method: "GET", Path: "/one", Response: models.Response{Status: 200}})
	assert.NoError(t, err)
	err = r.AddMock(&models.Mock{Name: "two", Method: "GET", Path: "/two", Response: models.Response{Status: 200}})
	assert.NoError(t, err)

	err = r.UpdateMock("one", &models.Mock{Name: "two", Method: "POST", Path: "/new", Response: models.Response{Status: 201}})
	assert.Error(t, err)

	_, err = r.GetMock("one")
	assert.NoError(t, err)
	_, err = r.GetMock("two")
	assert.NoError(t, err)
}
