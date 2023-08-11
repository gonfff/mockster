package parsers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gonfff/mockster/app/models"
	"gopkg.in/yaml.v2"
)

// ParseYAML reads a YAML file and returns a Mocks struct
func ParseYAML(filePath string) (*models.Mocks, error) {
	fileContent, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	mocks := &models.Mocks{}

	if err = yaml.Unmarshal(fileContent, mocks); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	if err = mocks.Validate(); err != nil {
		return nil, fmt.Errorf("error validating YAML: %w", err)
	}
	return mocks, nil
}
