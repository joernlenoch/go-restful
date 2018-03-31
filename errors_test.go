package restful_test

import (
  "testing"
  "github.com/joernlenoch/go-restful"
  "errors"
  "github.com/stretchr/testify/assert"
  "net/http"
)

func TestInvalidForm(t *testing.T) {
  t.Parallel()

  err := errors.New("test")
  resp := restful.InvalidForm(err)

  assert.NotNil(t, resp)
}

func TestInvalidJSON(t *testing.T) {
  t.Parallel()

  err := errors.New("test")
  resp := restful.InvalidJSON(err)

  assert.NotNil(t, resp)

}

func TestNotFound(t *testing.T) {
  t.Parallel()

  resp := restful.NotFound()
  assert.NotNil(t, resp)
}

func TestUnauthorized(t *testing.T) {
  t.Parallel()

  resp := restful.Unauthorized()

  assert.NotNil(t, resp)
}

func TestForbidden(t *testing.T) {
  t.Parallel()

  resp := restful.Forbidden()

  assert.NotNil(t, resp)
}

func TestBadRequest(t *testing.T) {
  t.Parallel()

  resp := restful.BadRequest("msg", "reason")

  assert.NotNil(t, resp)
}

func TestServerError(t *testing.T) {
  t.Parallel()

  err := errors.New("test")
  resp := restful.ServerError(err, "hello %s", "world")

  assert.Equal(t, resp.GetSource(), err, "should keep the source error")
  assert.Equal(t, len(resp.GetStack()), 1, "should have one stack entry")
  assert.Contains(t, resp.GetStack()[0], "hello world", "should have stored the stack info")
  assert.Equal(t, resp.GetMessage(), restful.ServerErrorText, "Must use the config variable")
  assert.Equal(t, resp.GetCode(), http.StatusInternalServerError, "must serve with 500")
}

func TestStack(t *testing.T) {
  t.Parallel()

  err := errors.New("test")
  resp := restful.Stack(err, "hello %s", "world")

  assert.Equal(t, resp.GetSource(), err, "should keep the source error")
  assert.Equal(t, len(resp.GetStack()), 1, "should have one stack entry")
  assert.Contains(t, resp.GetStack()[0], "hello world", "should have stored the stack info")
}