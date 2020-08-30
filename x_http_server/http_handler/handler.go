package http_handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/getto-systems/project-example-auth/request"
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
func (handler Handler) parseBodyProto(input protoreflect.ProtoMessage) (err error) {
	if handler.httpRequest.Body == nil {
		err = errEmptyBody
		return
	}

	body, err := ioutil.ReadAll(handler.httpRequest.Body)
	if err != nil {
		err = errBodyParseFailed
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		err = errBodyParseFailed
		return
	}

	err = proto.Unmarshal(decoded, input)
	if err != nil {
		return
	}

	return nil
}

func (handler Handler) ok(body interface{}) {
	handler.jsonResponse(http.StatusOK, body)
}

func (handler Handler) badRequest() {
	handler.jsonResponse(http.StatusBadRequest, newErrorResponseBody("bad-request"))
}
func (handler Handler) invalidTicket() {
	handler.jsonResponse(http.StatusUnauthorized, newErrorResponseBody("invalid-ticket"))
}
func (handler Handler) internalServerError() {
	handler.jsonResponse(http.StatusInternalServerError, newErrorResponseBody("internal-server-error"))
}

func (handler Handler) unauthorized(message string) {
	handler.jsonResponse(http.StatusUnauthorized, newErrorResponseBody(message))
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
