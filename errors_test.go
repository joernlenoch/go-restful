package restful_test

import (
	"encoding/json"
	"errors"
	"github.com/joernlenoch/go-restful"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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
	resp := restful.ServerError(err, "hello", "world")

	assert.Equal(t, resp.GetSource(), err, "should keep the source error")
	assert.Equal(t, len(resp.GetStack()), 1, "should have one stack entry")
	assert.Contains(t, resp.GetStack()[0], "hello, world", "should have stored the stack info")
	assert.Equal(t, resp.GetMessage(), "test", "Must set the errors text as default")
	assert.Equal(t, resp.GetCode(), http.StatusInternalServerError, "must serve with 500")
}

func TestStack(t *testing.T) {
	t.Parallel()

	err := errors.New("test")
	resp := restful.Stack(err, "hello", "world", 12)

	assert.Equal(t, resp.GetSource(), err, "should keep the source error")
	assert.Equal(t, 2, len(resp.GetStack()), "should have one stack entry")
	assert.Contains(t, resp.GetStack()[0], "hello, world, 12", "should have stored the stack info")
}

func TestStackf(t *testing.T) {
	t.Parallel()

	err := errors.New("test")
	resp := restful.Stackf(err, "hello %s %d", "world", 12)

	assert.Equal(t, resp.GetSource(), err, "should keep the source error")
	assert.Equal(t, 2, len(resp.GetStack()), "should have one stack entry")
	assert.Contains(t, resp.GetStack()[0], "hello world 12", "should have stored the stack info")
}

func TestDevelopment(t *testing.T) {
	restful.Development = false

	defer func() {
		restful.Development = true
	}()

	err := errors.New("test")
	resp := restful.ServerError(err, "test info")

	assert.Equal(t, resp.GetSource(), err, "should keep the source error")
	assert.Equal(t, len(resp.GetStack()), 1, "should have one stack entry")
	assert.Contains(t, resp.GetStack()[0], "test info", "should have stored the stack info")

	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "stack", "must not contain the stack information")
	assert.NotContains(t, string(data), "source", "must not contain the source information")
}

func TestResponse_SetMessage(t *testing.T) {
	err := errors.New("err")
	resp := restful.ServerError(err)

	resp.SetMessage("test")
	assert.Equal(t, resp.GetMessage(), "test", "should update the core message")
}
