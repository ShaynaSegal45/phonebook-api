package errors

import (
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	Err        error
	StatusCode int
	Wrapper    []string
}

const (
	InternalError = http.StatusInternalServerError
	ConflictError = http.StatusConflict
	NotFoundError = http.StatusNotFound
)

func CreateError(operationName, functionName string, err error, status ...int) *Error {
	if err != nil {
		code := InternalError
		if len(status) > 0 {
			code = status[0]
		}
		return &Error{
			Err:        err,
			StatusCode: code,
			Wrapper:    []string{fmt.Sprintf("%s.%s", operationName, functionName)},
		}
	}
	return nil
}

func (e *Error) Error() string {
	wrapper := strings.Join(e.Wrapper, " -> ")
	return fmt.Sprintf("error: %s, trace: %s", e.Err.Error(), wrapper)
}

func (e *Error) ErrorWrapper(operationName, functionName string) *Error {
	e.Wrapper = append(e.Wrapper, fmt.Sprintf("%s.%s", operationName, functionName))
	return e
}
