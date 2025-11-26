package domain

type ErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

type DuplicateError struct {
	Field string
}

func (e *DuplicateError) Error() string {
	return e.Field + " already exists"
}
