package model

import "fmt"

type missingField struct {
	Name string
}

func (e *missingField) Error() string {
	return fmt.Sprintf("Missing field: %s", e.Name)
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

func NewInvalidField(name string, reason string) *invalidField {
	return &invalidField{name, reason}
}
