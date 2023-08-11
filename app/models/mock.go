package models

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Mock is a struct for building endpoints for mocking
type Mock struct {
	Path     string   `yaml:"path" validate:"required,validHTTPPath"`
	Method   string   `yaml:"method" validate:"required,oneof=GET POST PUT PATCH DELETE OPTIONS HEAD"`
	Request  Request  `yaml:"request" validate:"dive"`
	Response Response `yaml:"response" validate:"dive"`
}

// Request is the request for getting Response
type Request struct {
	Headers     map[string]string `yaml:"headers"`
	QueryParams map[string]string `yaml:"query_params"`
	Cookies     map[string]string `yaml:"cookies"`
	Body        string            `yaml:"body"`
}

// Response is the response for Request
type Response struct {
	Status  int               `yaml:"status" validate:"required,min=100,max=599"`
	Headers map[string]string `yaml:"headers"`
	Cookies map[string]string `yaml:"cookies"`
	Body    string            `yaml:"body"`
}

// Mocks is a slice of Mock
type Mocks struct {
	Mocks []Mock `yaml:"mocks" validate:"dive"`
}

var validPathRegex = regexp.MustCompile(`^/.*$`)
var forbiddenCharsRegex = regexp.MustCompile(`[<>{}|\^\[\]"]`)

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

	err := validate.Struct(m.Mocks)
	if err != nil {
		return fmt.Errorf("error validating Mocks: %w", err)
	}
	return nil
}
