package restful

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"net/http"
	"runtime"
	"strings"
)

var (
	Development bool

	MsgServerError  = "server-error"
	MsgUnauthorized = "not-authenticated"
	MsgForbidden    = "access-denied"
	MsgNotFound     = "not-found"
)

// M is a simple string map for result parameters
type M map[string]interface{}

// Response
type Response interface {
	error

	// GetCode returns the code of this response
	GetCode() int

	// SetCode sets the response code
	SetCode(int)

	// GetTracking returns the tracking id
	GetTracking() string

	// GetReason set the reason for this
	GetReason() string

	// GetMessage returns the core message of this response
	GetMessage() string

	// SetMessage sets the core message
	SetMessage(string)

	// GetStack returns the custom error stack
	GetStack() []string

	// GetSource returns the original error
	GetSource() error

	push(string)
}

type response struct {
	Code int `json:"code, omitempty"`

	// Tracking number of an error
	Tracking string   `json:"tracking,omitempty"`
	Reason   string   `json:"reason,omitempty"`
	Message  string   `json:"message,omitempty"`
	Stack    []string `json:"stack,omitempty"`
	Source   error    `json:"source,omitempty"`
}

func (r response) Error() string {

	if Development {
		src := ""
		if r.Source != nil {
			src = r.Source.Error()
		}

		return fmt.Sprintf("(%d) %s [by %s] STACK %s (Source: %s)",
			r.Code,
			r.Message,
			r.Reason,
			strings.Join(r.Stack, "\n"),
			src,
		)
	}

	if len(r.Reason) > 0 {
		return fmt.Sprintf("%s (%s)", r.Message, r.Reason)
	}

	return fmt.Sprintf("%s", r.Message)
}

func (r response) MarshalJSON() ([]byte, error) {

	data := struct {
		Tracking string      `json:"tracking,omitempty"`
		Message  string      `json:"message"`
		Reason   string      `json:"reason,omitempty"`
		Stack    []string    `json:"stack,omitempty"`
		Source   interface{} `json:"source,omitempty"`
	}{
		Tracking: r.Tracking,
		Message:  r.Message,
		Reason:   r.Reason,
	}

	if Development {
		data.Stack = r.Stack

		if src, ok := r.Source.(Response); ok {
			data.Source = src
		} else if r.Source != nil {
			data.Source = r.Source.Error()
		}
	}

	return json.Marshal(data)
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

func (r *response) SetMessage(s string) {
	r.Message = s
}

func (r response) GetStack() []string {
	return r.Stack
}

func (r response) GetSource() error {
	return r.Source
}

func (r *response) push(s string) {
	if r.Stack == nil {
		r.Stack = []string{s}
	} else {
		r.Stack = append([]string{s}, r.Stack...)
	}
}

func (r *response) pop(s string) {
	if r.Stack == nil {
		r.Stack = []string{s}
	} else {
		r.Stack = append(r.Stack, s)
	}
}

func fromError(err error) Response {
	if err == nil {
		return newResponse()
	}

	if r, ok := err.(Response); ok {
		return r
	}

	return ServerError(err)
}

func Stack(err error, info ...interface{}) Response {

	// Create or restore the previous response structure
	r := fromError(err)

	fileInfo := printCallerInfo(1)
	entry := printStack(info...)

	if len(entry) > 0 {
		r.push(fmt.Sprintf("%s at %s", entry, fileInfo))
	} else {
		r.push(fileInfo)
	}

	return r
}

func Stackf(err error, info ...interface{}) Response {
	return Stack(err, printStackf(info...))
}

func printStack(info ...interface{}) string {
	if len(info) > 0 {
		parts := make([]string, len(info))
		for i, el := range info {
			parts[i] = fmt.Sprint(el)
		}

		return fmt.Sprint("{", strings.Join(parts, ", "), "}")
	}
	return ""
}

func printStackf(info ...interface{}) string {
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

func printCallerInfo(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("[%s:%d] ", file, line)
}

func newResponse(info ...interface{}) *response {
	return &response{
		Stack: []string{
			printCallerInfo(2) + printStack(info...),
		},
		Code:    http.StatusInternalServerError,
		Message: "",
		Source:  nil,
	}
}

func InvalidJSON(err error, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusBadRequest
	r.Message = "invalid-json"
	r.Source = err
	return r
}

func InvalidForm(err error, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusBadRequest
	r.Message = "invalid-form"
	r.Source = err
	return r
}

func Unauthorized(info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusUnauthorized
	r.Message = MsgUnauthorized
	return r
}

func UnauthorizedWithReason(reason string, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusUnauthorized
	r.Message = MsgUnauthorized
	r.Reason = reason
	return r
}

func Forbidden(info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusForbidden
	r.Message = MsgForbidden
	return r
}

func ForbiddenWithReason(reason string, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusForbidden
	r.Message = MsgForbidden
	r.Reason = reason
	return r
}

func NotFound(info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusNotFound
	r.Message = MsgNotFound
	return r
}

func NotFoundWithReason(reason string, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusNotFound
	r.Message = MsgNotFound
	r.Reason = reason

	return r
}

func BadRequest(msg string, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusBadRequest
	r.Message = msg
	r.Reason = ""
	return r
}

func BadRequestWithReason(msg string, reason string, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusBadRequest
	r.Message = msg
	r.Reason = reason
	return r
}

func ServerError(err error, info ...interface{}) Response {
	r := newResponse(info...)
	r.Code = http.StatusInternalServerError
	r.Message = err.Error()
	r.Source = err

	// Keep the tracing and stack information
	if prev, ok := err.(Response); ok {
		r.Tracking = prev.GetTracking()
		for _, s := range prev.GetStack() {
			r.pop(s)
		}
	}

	if len(r.Tracking) == 0 {
		r.Tracking = xid.New().String()
	}

	return r
}

func ServerErrorWithReason(err error, reason string, info ...interface{}) Response {
	r := ServerError(err, info...)
	baseResp := r.(*response)
	baseResp.Reason = reason
	return r
}
