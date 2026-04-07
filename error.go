package jsonapi

import "fmt"

func Error(err ErrorType, args ...any) APIError {
	return &_error{ErrorType: err, args: args}
}

type ErrorGroup interface {
	Split() []error
}

type APIError interface {
	error

	Status() int
	Code() int
	Title() string
	Details() string
}

type _error struct {
	ErrorType
	args []any
}

func (e *_error) Error() string {
	if e.details != "" {
		return e.Details()
	}

	return e.title
}

func (e *_error) Status() int {
	return e.status
}

func (e *_error) Code() int {
	return e.code
}

func (e *_error) Title() string {
	return e.title
}

func (e *_error) Details() string {
	return fmt.Sprintf(e.details, e.args...)
}
