package errors

import "fmt"

type HttpError interface {
	error
}

type httpError struct {
	code    string
	message string
}

func NewHttpError(code, message string) HttpError {
	return &httpError{
		code:    code,
		message: message,
	}
}

func (e httpError) Error() string {
	return fmt.Sprintf("request failed with status %s and error %s", e.code, e.message)
}
