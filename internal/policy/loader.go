package policy

import (
	"errors"
	"strings"
)

var ErrServiceRequired = errors.New("service name is required")

type Validator interface {
	Validate(string) error
}

type RequiredValidator struct{}

func (v *RequiredValidator) Validate(service string) error {
	if v == nil {
		return nil
	}
	if strings.TrimSpace(service) == "" {
		return ErrServiceRequired
	}
	return nil
}

type Config struct {
	Labels    map[string]string
	Validator Validator
}

func LoadDefault() *Config {
	var required *RequiredValidator
	return &Config{Validator: required}
}
