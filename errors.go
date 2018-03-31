package restful

import (
	"fmt"
	"net/http"
	"runtime"
	"github.com/rs/xid"
)

var (
	ServerErrorText = "server.error"
	UnauthorizedText = "not-authenticated"
	ForbiddenText = "access-denied"
	NotFoundText = "not-found"
)

type	M map[string]interface{}

//
//
//
/*type Error struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
	Stack   []string `json:"stack,omitempty"`
	Source  error `json:"source,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s caused %s (%s)", e.Reason, e.Message, e.Source)
}*/


//
//
//
type Error struct {
	Code int `json:"code, omitempty"`

	// Tracking number of an error
	Tracking string `json:"tracking,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
	Stack   []string `json:"stack,omitempty"`
	Source  error `json:"source,omitempty"`

	Errors []Error `json:"errors,omitempty"`
}

func (r Error) Error() string {
	return fmt.Sprintf("%d: %s (%s)", r.Code, r.Message, r.Reason)
}

func fromError(err error) *Error {
	if r, ok := err.(*Error); ok {
		return r
	}

	return ServerError(err)
}

func InvalidJSON(err error) error {
	return Error{
		Code:    400,
		Source:  err,
		Message: "invalid-json",
	}
}

func InvalidForm(err error) error {
	return Error{
		Code: 400,

		Stack:   []string{"Expected valid url form data"},
		Source:  err,
		Message: "The given url form data is invalid.",
		Reason:  "general",
	}
}

func Stack(err error, info ...interface{}) *Error {

	// Create or restore the previous response structure
	r := fromError(err)

	fileInfo := printCallerInfo()
	entry := printStack(info...)

	if len(entry) > 0 {
		if r.Stack == nil {
			r.Stack = []string{fileInfo + entry}
		} else {
			r.Stack = append(r.Stack, fileInfo+ entry)
		}
	}

	return r
}

func printStack(info ...interface{}) string {
	s := ""
	if len(info) > 0 {
		// Make sure that the first entry always is a string...
		s = fmt.Sprintf("%v", info[0])
	}
	if len(info) > 1 {
		s = fmt.Sprintf(s, info[1:]...)
	}
	return s
}

func printCallerInfo() string {
	_, fn, line, _ := runtime.Caller(2)
	return fmt.Sprintf("[%s:%d] ", fn, line)
}

func BadRequest(msg string, reason string, stack ...interface{}) *Error {
	return BadRequestEx(Error{
		Message: msg,
		Reason: reason,
		Stack: []string{printStack(stack...)},
	})
}

func BadRequestEx(err Error) *Error {
	return &Error{
		Code: http.StatusBadRequest,

		Message: err.Message,
		Reason:  err.Reason,
		Stack:   err.Stack,
		Source:  err.Source,
	}
}

func Unauthorized() *Error {
	return &Error{
		Code: http.StatusUnauthorized,
		Message: UnauthorizedText,
	}
}

func Forbidden() *Error {
	return &Error{
		Code: http.StatusForbidden,
		Message: ForbiddenText,
	}
}

func NotFound() *Error {
	return &Error{
		Code: http.StatusNotFound,
		Message: NotFoundText,
	}
}

func ServerError(err error, info ...interface{}) *Error {

	// To prevent wrong reuse with the old style...
	if _, ok := err.(*Error); ok {
		return Stack(err, info...)
	}

	var stack []string = nil
	if len(info) > 0 {
		stack = []string{printCallerInfo() +	printStack(info...)}
	}

	return &Error{
		Code: http.StatusInternalServerError,
		Message: ServerErrorText,
		Tracking: xid.New().String(),
		Stack: stack,
		Source: err,
	}
}
