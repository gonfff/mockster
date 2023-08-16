package parsers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gonfff/mockster/app/models"
	"gopkg.in/yaml.v2"
)

// ParseYAML parses YAML content and returns Mocks struct
func ParseYAML(content []byte) ([]models.Mock, error) {
	mocks := &models.Mocks{}

	if err := yaml.Unmarshal(content, mocks); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	if err := mocks.Validate(); err != nil {
		return nil, fmt.Errorf("error validating YAML: %w", err)
	}

	return mocks.Mocks, nil
}

// FileYAML parses YAML file from filesystem and returns mocks
func FileYAML(filePath string) ([]models.Mock, error) {
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return ParseYAML(content)
}

// ToYAML converts mocks to YAML
func ToYAML(mocks []*models.Mock) ([]byte, error) {
	m := &models.Mocks{Mocks: make([]models.Mock, 0, len(mocks))}
	for _, mock := range mocks {
		m.Mocks = append(m.Mocks, *mock)
	}
	data, err := yaml.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error marshaling mocks to YAML: %w", err)
	}

	return data, nil
}
