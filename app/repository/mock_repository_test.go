package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitRepository(t *testing.T) {
	repo, err := InitRepository()
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}
