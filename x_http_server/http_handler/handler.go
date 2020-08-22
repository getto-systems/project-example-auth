package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/request"
)

var (
	errEmptyBody       = errors.New("empty body")
	errBodyParseFailed = errors.New("body parse failed")
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
	return request.NewRequest(request.RequestedAtNow(), request.RemoteAddr(handler.httpRequest.RemoteAddr))
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
	if password.ErrCheck.IsSameCategory(err) {
		handler.invalidPassword()
		return
	}
	if credential.ErrClearCredential.IsSameCategory(err) {
		handler.invalidCredential()
		return
	}

	switch err {
	case errEmptyBody,
		errBodyParseFailed,
		errTicketTokenNotFound:

		handler.badRequest()
		return
	}

	handler.internalServerError()
}

func (handler Handler) ok(body interface{}) {
	handler.jsonResponse(http.StatusOK, body)
}

func (handler Handler) invalidPassword() {
	handler.jsonResponse(http.StatusUnauthorized, newErrorResponseBody("invalid-password"))
}
func (handler Handler) invalidCredential() {
	handler.jsonResponse(http.StatusUnauthorized, newErrorResponseBody("invalid-credential"))
}
func (handler Handler) badRequest() {
	handler.jsonResponse(http.StatusBadRequest, newErrorResponseBody("bad-request"))
}
func (handler Handler) internalServerError() {
	handler.jsonResponse(http.StatusInternalServerError, newErrorResponseBody("internal-server-error"))
}

func newErrorResponseBody(message string) errorResponseBody {
	return errorResponseBody{
		Message: message,
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
