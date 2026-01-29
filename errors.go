package ngrokd

import (
	"errors"
	"fmt"
)

var (
	ErrEndpointNotFound = errors.New("endpoint not found")
	ErrClosed           = errors.New("dialer is closed")
)

// EndpointNotFoundError is returned when an endpoint is not in cache
type EndpointNotFoundError struct {
	Hostname string
}

func (e *EndpointNotFoundError) Error() string {
	return fmt.Sprintf("endpoint not found: %s", e.Hostname)
}

func (e *EndpointNotFoundError) Is(target error) bool {
	return target == ErrEndpointNotFound
}
