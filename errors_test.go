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
}

func TestInvalidJSON(t *testing.T) {
  t.Parallel()
}

func TestNotFound(t *testing.T) {
  t.Parallel()
}

func TestUnauthorized(t *testing.T) {
  t.Parallel()
}

func TestForbidden(t *testing.T) {
  t.Parallel()
}

func TestBadRequest(t *testing.T) {
  t.Parallel()
}

func TestBadRequestEx(t *testing.T) {
  t.Parallel()
}

func TestServerError(t *testing.T) {
  t.Parallel()

  err := errors.New("test")
  se := restful.ServerError(err, "hello %s", "world")

  assert.Equal(t, se.GetSource(), err, "should keep the source error")
  assert.Equal(t, len(se.GetStack()), 1, "should have one stack entry")
  assert.Contains(t, se.GetStack()[0], "hello world", "should have stored the stack info")
  assert.Equal(t, se.GetMessage(), restful.ServerErrorText, "Must use the config variable")
  assert.Equal(t, se.GetCode(), http.StatusInternalServerError, "must serve with 500")
}

func TestStack(t *testing.T) {
  t.Parallel()

  err := errors.New("test")
  stacked := restful.Stack(err, "hello %s", "world")

  assert.Equal(t, stacked.GetSource(), err, "should keep the source error")
  assert.Equal(t, len(stacked.GetStack()), 1, "should have one stack entry")
  assert.Contains(t, stacked.GetStack()[0], "hello world", "should have stored the stack info")
}