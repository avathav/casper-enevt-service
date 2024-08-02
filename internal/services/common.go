package services

import "fmt"

type ErrResourceIsRequired struct {
	service  string
	resource string
}

func NewErrResourceIsRequired(service string, resource string) *ErrResourceIsRequired {
	return &ErrResourceIsRequired{
		service:  service,
		resource: resource,
	}
}

func (e ErrResourceIsRequired) Error() string {
	return fmt.Sprintf("resource: %s is required for %s service", e.resource, e.service)
}
