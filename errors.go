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
/*type error struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
	Stack   []string `json:"stack,omitempty"`
	Source  error `json:"source,omitempty"`
}

func (e error) error() string {
	return fmt.Sprintf("%s caused %s (%s)", e.Reason, e.Message, e.Source)
}*/

type Response interface {
	error

  GetCode() int
  SetCode(int)

  GetTracking() string

  GetReason() string

  GetMessage() string
  SetMessage(string)

  GetStack() []string
  Push(s string)

  GetSource() error
}


//
//
//
type response struct {
	Code int `json:"code, omitempty"`

	// Tracking number of an error
	Tracking string `json:"tracking,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
	Stack   []string `json:"stack,omitempty"`
	Source  error `json:"source,omitempty"`
}

func (r response) Error() string {
  return fmt.Sprintf("%d: %s (%s)", r.Code, r.Message, r.Reason)
}

func (r response) GetCode() int {
  return r.Code
}

func (r *response) SetCode(c int) {
	r.Code = c
}

func (r response) GetTracking() string {
  return r.Tracking
}

func (r response) GetReason() string {
  return r.Reason
}

func (r response) GetMessage() string {
  return r.Message
}

func (r response) SetMessage(s string) {
	r.Message = s
}

func (r response) GetStack() []string {
  return r.Stack
}

func (r response) GetSource() error {
  return r.Source
}

func (r *response) Push(s string) {
  if r.Stack == nil {
    r.Stack = []string{s}
  } else {
    r.Stack = append(r.Stack, s)
  }
}


func fromError(err error) Response {
	if r, ok := err.(Response); ok {
		return r
	}

	return ServerError(err)
}

func InvalidJSON(err error) Response {
	return &response{
		Code:    400,
		Source:  err,
		Message: "invalid-json",
	}
}

func InvalidForm(err error) Response {
	return &response{
		Code: 400,

		Stack:   []string{"Expected valid url form data"},
		Source:  err,
		Message: "The given url form data is invalid.",
		Reason:  "general",
	}
}

func Stack(err error, info ...interface{}) Response {

	// Create or restore the previous response structure
	r := fromError(err)

	fileInfo := printCallerInfo()
	entry := printStack(info...)

	if len(entry) > 0 {
    r.Push(fileInfo + entry)
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

func BadRequest(msg string, reason string, info ...interface{}) Response {

  var stack []string = nil
  if len(info) > 0 {
    stack = []string{printStack(info...)}
  }

	return &response{
		Message: msg,
		Reason: reason,
		Stack: stack,
	}
}

func Unauthorized() Response {
	return &response{
		Code: http.StatusUnauthorized,
		Message: UnauthorizedText,
	}
}

func Forbidden() Response {
	return &response{
		Code: http.StatusForbidden,
		Message: ForbiddenText,
	}
}

func NotFound() Response {
	return &response{
		Code: http.StatusNotFound,
		Message: NotFoundText,
	}
}

func ServerError(err error, info ...interface{}) Response {

	// To prevent wrong reuse with the old style...
	if _, ok := err.(Response); ok {
		return Stack(err, info...)
	}

	var stack []string = nil
	if len(info) > 0 {
		stack = []string{printCallerInfo() +	printStack(info...)}
	}

	return &response{
		Code: http.StatusInternalServerError,
		Message: ServerErrorText,
		Tracking: xid.New().String(),
		Stack: stack,
		Source: err,
	}
}
