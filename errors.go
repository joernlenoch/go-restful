package restful

import (
	"fmt"
	"net/http"
	"runtime"
)

type (
	Error struct {
		error `json:"-"`

		Reason  string `json:"reason"`
		Message string `json:"message"`
		DevInfo string `json:"devInfo,omitempty"`
		Source  string `json:"source,omitempty"`
	}

	Response struct {
		error `json:"-"`

		Code int `json:"code"`

		Tracking string `json:"-"`
		Reason   string `json:"reason,omitempty"`
		Message  string `json:"message"`
		DevInfo  string `json:"devInfo,omitempty"`
		Source   string `json:"source,omitempty"`

		Errors []Error `json:"errors,omitempty"`
	}

	M map[string]interface{}
)

func InvalidJSON(err error) Response {
	return Response{
		Code:    400,
		Source:  err.Error(),
		Message: "invalid-json",
	}
}

func InvalidForm(err error) Response {
	return Response{
		Code: 400,

		DevInfo: "Expected valid url form data",
		Source:  err.Error(),
		Message: "The given url form data is invalid.",
		Reason:  "general",
	}
}

func (e Response) Error() string {
	return fmt.Sprintf("%d: %s (%s)", e.Code, e.Message, e.Reason)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s caused %s (%s)", e.Reason, e.Message, e.Source)
}

/*
func Send(ctx *context.ExtendedContext, err error) error {

	var resp Response

	resp, ok := err.(Response)
	if !ok {
		resp = ServerError(resp, "Unknown error")
	}

	// Reset all development information in production mode.
	if !Development {
		resp.DevInfo = ""
		resp.Source = ""
	}

	// Translate the message as good as possible
	resp.Message = i18n.Translate(ctx, resp.Message)

	switch {
	case 200 <= resp.Code && resp.Code < 400 || resp.Code == 404:
		return ctx.JSON(iris.StatusOK, resp)
	case resp.Code == iris.StatusUnauthorized:
		return ctx.JSON(iris.StatusUnauthorized, resp)
	case resp.Code >= 400 && resp.Code < 500:
		return ctx.JSON(iris.StatusBadRequest, resp)
	case resp.Code >= 500:
		log.Printf("Internal Server Error %#v", resp)
		return ctx.JSON(iris.StatusInternalServerError, resp)
	}

	log.Panicf("An unsupported status code has been used: %#v", resp)
	return nil
}


func OK(ctx *context.ExtendedContext, v interface{}) error {
	return ctx.JSON(iris.StatusOK, v)
}
*/

func BadRequest(err Error) Response {
	return Response{
		Code: http.StatusBadRequest,

		Message: err.Message,
		Reason:  err.Reason,
		DevInfo: err.DevInfo,
		Source:  err.Source,
	}
}

func Unauthorized() Response {
	return Response{
		Code: http.StatusUnauthorized,

		Message: "not-authenticated",
		// Reason:  err.Reason,
		// DevInfo: err.DevInfo,
		// Source:  err.Source,
	}
}

func Forbidden() Response {
	return Response{
		Code: http.StatusForbidden,

		Message: "no-access",
		// Reason:  err.Reason,
		// DevInfo: err.DevInfo,
		// Source:  err.Source,
	}
}

func NotFound() Response {
	return Response{
		Code: http.StatusNotFound,

		Message: "not-found",
		// Reason:  err.Reason,
		// DevInfo: err.DevInfo,
		// Source:  err.Source,
	}
}

func ServerError(err error, comments ...interface{}) Response {

	_, fn, line, _ := runtime.Caller(1)

	// var tracking string
	response := Response{
		Code: http.StatusInternalServerError,
	}

	if len(comments) > 1 {
		response.DevInfo = fmt.Sprintf("%v [%s:%d] %v", comments[0], fn, line, comments[1:])
	} else if len(comments) > 0 {
		response.DevInfo = fmt.Sprintf("%v [%s:%d]", comments[0], fn, line)
	}

	// Check if we have a tracing problem here.
	if oldResp, ok := err.(Response); ok {
		// tracking = oldResp.Tracking
		response.Source = oldResp.Source
		response.DevInfo = fmt.Sprintf("%s ==> %s", response.DevInfo, oldResp.DevInfo)
	} else {
		// tracking = gocql.UUIDFromTime(time.Now()).String()

		if oldError, ok := err.(Error); ok {
			response.Source = oldError.Source

			if len(oldError.DevInfo) > 0 {
				response.DevInfo = fmt.Sprintf("%s ==> %s", response.DevInfo, oldError.DevInfo)
			}

			if len(oldError.Message) > 0 {
				response.DevInfo = fmt.Sprintf("%s ==> %s", response.DevInfo, oldError.Message)
			}

		} else if err != nil {
			response.Source = err.Error()
		}

	}

	// response.Tracking = tracking
	response.Message = fmt.Sprintf(
		"An unexpected Error occured. If this Error persists. Please contact the support.",
		// tracking,
	)

	return response
}
