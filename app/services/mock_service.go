package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gonfff/mockster/app/models"
	"github.com/gonfff/mockster/app/repository"
)

// MockService orchestrates mock operations.
type MockService struct {
	repo repository.MockRepository
}

// NewMockService creates a new mock service.
func NewMockService(repo repository.MockRepository) *MockService {
	return &MockService{repo: repo}
}

// MatchMock returns a matching mock for request details.
func (s *MockService) MatchMock(method, path, body string, headers http.Header, cookies []*http.Cookie, query url.Values) (*models.Mock, error) {
	endpoint := fmt.Sprintf("%v %v", method, path)
	mockNames, err := s.repo.GetMockNames(endpoint)
	if err != nil {
		return nil, err
	}
	if len(mockNames) == 0 {
		return nil, fmt.Errorf("no mocks found for endpoint %s", endpoint)
	}

	body = sanitizeBody(body)

	var matchErr error
	for _, mockName := range mockNames {
		mock, getErr := s.repo.GetMock(mockName)
		if getErr != nil {
			matchErr = getErr
			continue
		}
		if validateRequest(mock, headers, cookies, query, body) == nil {
			return mock, nil
		}
		matchErr = fmt.Errorf("no matching mock variant found")
	}

	if matchErr == nil {
		matchErr = fmt.Errorf("no matching mock variant found")
	}
	return nil, matchErr
}

func (s *MockService) GetMocks() ([]*models.Mock, error) {
	return s.repo.GetMocks()
}

func (s *MockService) CreateMock(mock *models.Mock) error {
	return s.repo.AddMock(mock)
}

func (s *MockService) DeleteMock(name string) error {
	return s.repo.DeleteMock(name)
}

func (s *MockService) UpdateMock(name string, mock *models.Mock) error {
	return s.repo.UpdateMock(name, mock)
}

func (s *MockService) ReplaceAll(mocks []models.Mock) error {
	items := make([]*models.Mock, 0, len(mocks))
	for i := range mocks {
		m := mocks[i]
		items = append(items, &m)
	}
	return s.repo.ReplaceAll(items)
}

func validateRequest(mock *models.Mock, headers http.Header, cookies []*http.Cookie, query url.Values, body string) error {
	for k, v := range mock.Request.Headers {
		if headers.Get(k) != v {
			return fmt.Errorf("header %v is not equal to %v", k, v)
		}
	}

	cookieMap := make(map[string]string, len(cookies))
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}
	for k, v := range mock.Request.Cookies {
		if cookieMap[k] != v {
			return fmt.Errorf("cookie %v is not equal to %v", k, v)
		}
	}

	for k, v := range mock.Request.QueryParams {
		if query.Get(k) != v {
			return fmt.Errorf("query param %v is not equal to %v", k, v)
		}
	}

	if mock.Request.Body != body {
		return fmt.Errorf("body is not equal to %v", mock.Request.Body)
	}

	return nil
}

func sanitizeBody(body string) string {
	body = strings.ReplaceAll(body, "\n", "")
	body = strings.ReplaceAll(body, "\t", "")
	return body
}
