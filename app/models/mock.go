package models

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Mock is a struct for building endpoints for mocking
type Mock struct {
	Name     string   `yaml:"name" json:"name" validate:"required"`
	Path     string   `yaml:"path" json:"path" validate:"required,validHTTPPath"`
	Method   string   `yaml:"method" json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE OPTIONS HEAD"`
	Request  Request  `yaml:"request" json:"request" validate:"dive"`
	Response Response `yaml:"response" json:"response" validate:"required,dive"`
}

// Request is the request for getting Response
type Request struct {
	Headers     map[string]string `yaml:"headers" json:"headers"`
	QueryParams map[string]string `yaml:"query_params" json:"query_params"`
	Cookies     map[string]string `yaml:"cookies" json:"cookies"`
	Body        string            `yaml:"body" json:"body"`
}

// Response is the response for Request
type Response struct {
	Status  int               `yaml:"status" json:"status" validate:"required,min=100,max=599"`
	Headers map[string]string `yaml:"headers" json:"headers"`
	Cookies map[string]string `yaml:"cookies" json:"cookies"`
	Body    string            `yaml:"body" json:"body"`
}

// Mocks is a slice of Mock
type Mocks struct {
	Mocks []Mock `yaml:"mocks" validate:"dive"`
}

var validPathRegex = regexp.MustCompile(`^/.*$`)
var forbiddenCharsRegex = regexp.MustCompile(`[<>{}|\^\[\]"]?&`)

func validHTTPPath(fl validator.FieldLevel) bool {
	path := fl.Field().String()

	// Проверка на валидный HTTP путь
	if !validPathRegex.MatchString(path) {
		return false
	}

	// Проверка на запрещенные символы
	if forbiddenCharsRegex.MatchString(path) {
		return false
	}

	return true
}

// Validate validates the Mocks struct
func (m *Mocks) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("validHTTPPath", validHTTPPath)

	err := validate.Struct(m)
	if err != nil {
		return fmt.Errorf("error validating Mocks: %w", err)
	}
	return nil
}

// Validate validates the Mock struct
func (m *Mock) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("validHTTPPath", validHTTPPath)

	err := validate.Struct(m)
	if err != nil {
		return fmt.Errorf("error validating Mock: %w", err)
	}
	return nil
}
