package services

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gonfff/mockster/app/models"
	"github.com/gonfff/mockster/app/repository"
	"github.com/stretchr/testify/assert"
)

func TestMatchMockMatchesSecondCandidate(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	_ = repo.AddMock(&models.Mock{
		Name:   "first",
		Method: "POST",
		Path:   "/test",
		Request: models.Request{
			Body: "one",
		},
		Response: models.Response{Status: 200, Body: "first"},
	})
	_ = repo.AddMock(&models.Mock{
		Name:   "second",
		Method: "POST",
		Path:   "/test",
		Request: models.Request{
			Body: "two",
		},
		Response: models.Response{Status: 200, Body: "second"},
	})

	svc := NewMockService(repo)
	m, err := svc.MatchMock(http.MethodPost, "/test", "two", http.Header{}, []*http.Cookie{}, url.Values{})
	assert.NoError(t, err)
	assert.Equal(t, "second", m.Name)
}
