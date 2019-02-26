package linter

type (
	// Error represents a linting with a resource
	// so it can be formatted separately with different
	// colours when printing to user.
	Error struct {
		error
		Resource string
	}
)

// NewError creates a new error.
func NewError(err error, resource string) Error {
	return Error{
		error:    err,
		Resource: resource,
	}
}

// Cause returns the underlying error this error
// is embedding.
func (e Error) Cause() error {
	return e.error
}
