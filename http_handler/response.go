package http_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type Response struct {
	cookie Cookie
}

type authenticatedResponseBody struct {
	Token []byte `json:"token"`
	Nonce string `json:"nonce"`
}

type errorResponseBody struct {
	Message string `json:"message"`
}

func NewResponse(cookie Cookie) Response {
	return Response{
		cookie: cookie,
	}
}

func (response Response) ResetCookie(w http.ResponseWriter) {
	response.cookie.resetTicket(w)
	response.cookie.resetContentToken(w)
}

func (response Response) Authenticated(
	w http.ResponseWriter,
	ticket ticket.Ticket,
	nonce ticket.Nonce,
	apiToken ticket.ApiToken,
	contentToken ticket.ContentToken,
	expires data.Expires,
	logger Logger,
) {
	response.cookie.setTicket(w, ticket, expires)
	response.cookie.setContentToken(w, contentToken, expires)

	JsonResponse(w, http.StatusOK, authenticatedResponseBody{
		Token: apiToken,
		Nonce: string(nonce),
	})
}

func (response Response) OK(w http.ResponseWriter) {
	JsonResponse(w, http.StatusOK, "OK")
}

func (response Response) Error(w http.ResponseWriter, err error) {
	JsonResponse(w, status(err), errorResponseBody{
		Message: err.Error(),
	})
}

func (response Response) BadRequest(w http.ResponseWriter, err error) {
	JsonResponse(w, http.StatusBadRequest, errorResponseBody{
		Message: err.Error(),
	})
}

func (response Response) Unauthorized(w http.ResponseWriter, err error) {
	JsonResponse(w, http.StatusUnauthorized, errorResponseBody{
		Message: err.Error(),
	})
}

func (response Response) InternalServerError(w http.ResponseWriter, err error) {
	// TODO ここで slack とかにログを送信したい
	JsonResponse(w, http.StatusInternalServerError, errorResponseBody{
		Message: err.Error(),
	})
}

func status(err error) int {
	switch err {
	case
		ErrEmptyBody,
		ErrBodyParseFailed:

		return http.StatusBadRequest

	case
		ErrTicketCookieNotFound,
		ticket.ErrValidateFailed,
		ticket.ErrExtendFailed:

		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}

func JsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", jsonData)
}
