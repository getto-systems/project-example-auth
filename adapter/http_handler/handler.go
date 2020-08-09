package http_handler

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/errors"
)

var (
	errEmptyBody       = goerrors.New("empty body")
	errBodyParseFailed = goerrors.New("body parse failed")
)

type (
	Handler struct {
		httpResponseWriter http.ResponseWriter
		httpRequest        *http.Request
	}

	errorResponseBody struct {
		Message string `json:"message"`
	}
)

func NewHandler(w http.ResponseWriter, r *http.Request) Handler {
	return Handler{
		httpResponseWriter: w,
		httpRequest:        r,
	}
}

func (handler Handler) newRequest() request.Request {
	return request.NewRequest(time.Now(), request.RemoteAddr(handler.httpRequest.RemoteAddr))
}

func (handler Handler) parseBody(input interface{}) (err error) {
	if handler.httpRequest.Body == nil {
		err = errEmptyBody
		return
	}

	err = json.NewDecoder(handler.httpRequest.Body).Decode(&input)
	if err != nil {
		err = errBodyParseFailed
		return
	}

	return nil
}

func (handler Handler) errorResponse(err error) {
	if goerrors.Is(err, errors.ErrPasswordCheck) {
		handler.badRequest(err)
		return
	}
	if goerrors.Is(err, errors.ErrTicketValidate) {
		handler.unauthorized(err)
		return
	}

	switch err {
	case errEmptyBody,
		errBodyParseFailed:

		handler.badRequest(err)
		return

	case errTicketTokenNotFound:
		handler.unauthorized(err)
		return
	}

	handler.internalServerError(err)
}

func (handler Handler) ok(body interface{}) {
	handler.jsonResponse(http.StatusOK, body)
}
func (handler Handler) badRequest(err error) {
	handler.jsonResponse(http.StatusBadRequest, newErrorResponseBody(err))
}
func (handler Handler) unauthorized(err error) {
	handler.jsonResponse(http.StatusUnauthorized, newErrorResponseBody(err))
}
func (handler Handler) internalServerError(err error) {
	handler.jsonResponse(http.StatusInternalServerError, newErrorResponseBody(err))
}

func newErrorResponseBody(err error) errorResponseBody {
	return errorResponseBody{
		Message: err.Error(),
	}
}

func (handler Handler) jsonResponse(status int, data interface{}) {
	handler.httpResponseWriter.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		handler.httpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.httpResponseWriter.WriteHeader(status)
	fmt.Fprintf(handler.httpResponseWriter, "%s", jsonData)
}
