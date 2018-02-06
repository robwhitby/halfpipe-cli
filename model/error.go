package model

import (
	"fmt"
)

type Documented interface {
	DocumentationPath() string
}

type missingField struct {
	Name string
}

func (e *missingField) Error() string {
	return fmt.Sprintf("Missing field: %s", e.Name)
}

func (e *missingField) DocumentationPath() string {
	return "/docs/manifest/fields#" + e.Name
}

func NewMissingField(name string) *missingField {
	return &missingField{name}
}

type invalidField struct {
	Name   string
	Reason string
}

func (e *invalidField) Error() string {
	return fmt.Sprintf("Invalid value for '%s': %s", e.Name, e.Reason)
}

func (e *invalidField) DocumentationPath() string {
	return "/docs/manifest/fields#" + e.Name
}

func NewInvalidField(name string, reason string) *invalidField {
	return &invalidField{name, reason}
}
