package status

import "net/http"

// Checker performs status checks for HTTP services.
type Checker interface {
	Check(Service) error
}

// Service represents a service which status can be checked.
type Service interface {
	NewStatusRequest() (*http.Request, error)
}
