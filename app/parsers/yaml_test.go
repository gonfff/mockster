package parsers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gonfff/mockster/app/models"
	"github.com/stretchr/testify/assert"
)

type yamlTestData struct {
	content []byte
	errMsg  string
}

var (
	validYAML = yamlTestData{content: []byte(`
mocks:
  - name: ping
    path: /ping
    method: GET
    response:
      status: 200
  - name: pong
    path: /pong
    method: GET
    response:
      status: 200
`), errMsg: ""}
	invalidYAML = yamlTestData{content: []byte(`
mocks:
  - name: ping
    path: /ping
    method: GET
    response:
      status: 200
  - name: pong
	path: /pong
    method: GET
    response:
      status: 200
`), errMsg: "error unmarshaling YAML:"}

	notValidYAML = yamlTestData{
		content: []byte(`
mocks:
  - name: ping
    path: /ping
    method: GET
    response:
      status: 200
  - name: pong
    path: /pong
    method: GET
    response:
      status: 900
`), errMsg: "error validating YAML:"}
)

func TestParseYAML(t *testing.T) {
	testCases := []struct {
		name string
		data yamlTestData
	}{
		{
			name: "Valid YAML",
			data: validYAML,
		},
		{
			name: "Invalid YAML",
			data: invalidYAML,
		},
		{
			name: "Not Valid YAML",
			data: notValidYAML,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mocks, err := ParseYAML(tc.data.content)
			if tc.data.errMsg == "" {
				assert.NoError(t, err)
				assert.Len(t, mocks, 2)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.data.errMsg)
			}
		})
	}
}

func TestFileYAML(t *testing.T) {
	testCases := []struct {
		name     string
		data     yamlTestData
		fileName string
	}{
		{
			name:     "Valid YAML",
			data:     validYAML,
			fileName: "test.yaml",
		},
		{
			name:     "Invalid YAML",
			data:     invalidYAML,
			fileName: "test.yaml",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tmpFile := createTempYAMLFile(t, testCase.data.content, testCase.fileName)
			defer os.Remove(tmpFile)

			if testCase.data.errMsg == "" {
				mocks, err := FileYAML(tmpFile)
				assert.NoError(t, err)
				assert.Len(t, mocks, 2)
			} else {
				_, err := FileYAML(tmpFile)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.data.errMsg)
			}

		})
	}
	t.Run("File not exists", func(t *testing.T) {
		_, err := FileYAML("non-existing-file.yaml")
		assert.Error(t, err)
	})
}

// createTempYAMLFile helper for creating temporary YAML file
func createTempYAMLFile(t *testing.T, content []byte, fileName string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, fileName)
	err := os.WriteFile(tmpFile, content, 0o0644)
	assert.NoError(t, err)
	return tmpFile
}

func TestToYAML(t *testing.T) {
	mock1 := &models.Mock{Name: "ping", Path: "/ping", Method: "GET", Response: models.Response{Status: 200}}
	mock2 := &models.Mock{Name: "pong", Path: "/pong", Method: "GET", Response: models.Response{Status: 200}}

	mocks := []*models.Mock{mock1, mock2}

	yamlData, err := ToYAML(mocks)

	assert.NoError(t, err)
	assert.NotEmpty(t, yamlData)
}
