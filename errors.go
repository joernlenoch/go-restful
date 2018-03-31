package restful

import (
  "fmt"
  "net/http"
  "runtime"
  "github.com/rs/xid"
  "strings"
  "encoding/json"
)

var (
  Development bool

  MsgServerError  = "server.error"
  MsgUnauthorized = "not-authenticated"
  MsgForbidden    = "access-denied"
  MsgNotFound     = "not-found"
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
  push(string)

  GetStack() []string

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

  if Development {
    src := ""
    if r.Source != nil {
      src = r.Source.Error()
    }

    return fmt.Sprintf("CODE %d MSG %s REASON %s STACK %s SOURCE %s",
      r.Code,
      r.Message,
      r.Reason,
      strings.Join(r.Stack, ","),
      src,
    )
  }

  if len(r.Reason) > 0 {
    return fmt.Sprintf("%s (%s)", r.Message, r.Reason)
  }

  return fmt.Sprintf("%s", r.Message)
}

func (r response) MarshalJSON() ([]byte, error) {

  data := struct{
    Tracking string `json:"tracking,omitempty"`
    Message string `json:"message"`
    Reason string `json:"reason,omitempty"`
    Stack []string `json:"stack,omitempty"`
    Source error `json:"source,omitempty"`
  }{
    Tracking: r.Tracking,
    Message: r.Message,
    Reason: r.Reason,
  }

  if Development {
    data.Stack = r.Stack
    data.Source = r.Source
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

func (r response) SetMessage(s string) {
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
  if r, ok := err.(Response); ok {
    return r
  }

  return ServerError(err)
}

func Stack(err error, info ...interface{}) Response {

  // Create or restore the previous response structure
  r := fromError(err)

  fileInfo := printCallerInfo(2)
  entry := printStack(info...)

  if len(entry) > 0 {
    r.push(fmt.Sprintf("%s at %s", entry, fileInfo))
  } else {
    r.push(fileInfo)
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

func printCallerInfo(skip int) string {
  _, file, line, _ := runtime.Caller(skip)
  return fmt.Sprintf("[%s:%d] ", file, line)
}

func newResponse(info ...interface{}) *response {
  return &response{
    Stack: []string{
      printCallerInfo(3) + printStack(info...),
    },
  }
}

func InvalidJSON(err error) Response {
  r := newResponse()
  r.Code = http.StatusBadRequest
  r.Message = "invalid-json"
  r.Source = err
  return r
}

func InvalidForm(err error) Response {
  r := newResponse()
  r.Code = http.StatusBadRequest
  r.Message = "invalid-form"
  r.Source = err
  return r
}

func Unauthorized() Response {
  r := newResponse()
  r.Code = http.StatusUnauthorized
  r.Message = MsgUnauthorized
  return r
}

func Forbidden() Response {
  r := newResponse()
  r.Code = http.StatusForbidden
  r.Message = MsgForbidden
  return r
}

func NotFound() Response {
  r := newResponse()
  r.Code = http.StatusNotFound
  r.Message = MsgNotFound
  return r
}

func BadRequest(msg string, reason string, info ...interface{}) Response {
  r := newResponse(info...)
  r.Code = http.StatusBadRequest
  r.Message = msg
  r.Reason = reason
  return r
}

func ServerError(err error, info ...interface{}) Response {

  r := newResponse(info...)
  r.Code = http.StatusInternalServerError
  r.Message = MsgServerError
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
