package errors

import (
	"fmt"
	"net/http"
	"strings"
)

// Op describes an operation
type Op string

// Custom error type
type Error struct {
	Err     error  `json:"-"`
	Status  int    `json:"httpStatus"`
	Op      Op     `json:"-"`
	Message string `json:"message"`
}

// E creates a new error instance. The extras interface can be an error, a message for the client or a HTTP status
func E(op Op, extras ...interface{}) error {
	e := &Error{
		Op:      op,
		Status:  http.StatusInternalServerError,
		Message: "Oops, something blew up!",
	}

	for _, ex := range extras {
		switch t := ex.(type) {
		case int:
			e.Status = t
		case error:
			e.Err = t
		case string:
			e.Message = t
		}
	}

	return e
}

// Error returns a string with information about the error for debugging.
// This value should not be returned to the user.
func (e *Error) Error() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf("%s: ", string(e.Op)))

	if e.Err != nil {
		b.WriteString(e.Err.Error())
	}

	return b.String()
}
