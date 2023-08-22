package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validMock           = Mock{Name: "ping", Path: "/ping", Method: "GET", Request: Request{}, Response: Response{Status: 200}}
	wrongMock           = Mock{Name: "pong", Path: "pong", Method: "GET", Request: Request{}, Response: Response{Status: 200}}
	wrongMockBadSymbols = Mock{Name: "/[]po|ng", Path: "pong", Method: "GET", Request: Request{}, Response: Response{Status: 200}}
	validMocks          = Mocks{Mocks: []Mock{validMock}}
	wrongMocks          = Mocks{Mocks: []Mock{wrongMock}}
)

func TestMocks_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		mocks Mocks
		err   bool
	}{
		{name: "valid", mocks: validMocks, err: false},
		{name: "wrong", mocks: wrongMocks, err: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mocks.Validate()
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMock_Validate(t *testing.T) {
	testCases := []struct {
		name string
		mock Mock
		err  bool
	}{
		{name: "valid", mock: validMock, err: false},
		{name: "wrong", mock: wrongMock, err: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mock.Validate()
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMock_ValidatePath(t *testing.T) {
	testCases := []struct {
		name string
		mock Mock
		err  bool
	}{
		{name: "valid", mock: validMock, err: false},
		{name: "wrong", mock: wrongMock, err: true},
		{name: "wrong bad symbols", mock: wrongMockBadSymbols, err: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mock.Validate()
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
