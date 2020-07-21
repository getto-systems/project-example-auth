package http_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
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

type getLoginResponseBody struct {
	LoginID string `json:"login_id"`
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

	jsonResponse(w, http.StatusOK, authenticatedResponseBody{
		Token: apiToken,
		Nonce: string(nonce),
	})
}

func (response Response) OK(w http.ResponseWriter) {
	jsonResponse(w, http.StatusOK, "OK")
}

func (response Response) Login(w http.ResponseWriter, login password.Login) {
	jsonResponse(w, http.StatusOK, getLoginResponseBody{
		LoginID: string(login.ID()),
	})
}

func (response Response) Error(w http.ResponseWriter, err error) {
	jsonResponse(w, status(err), errorResponseBody{
		Message: err.Error(),
	})
}

func status(err error) int {
	switch err {
	case
		ErrEmptyBody,
		ErrBodyParseFailed,
		password.ErrRegisterFailed:

		return http.StatusBadRequest

	case
		ErrTicketCookieNotFound,
		ticket.ErrValidateFailed,
		ticket.ErrExtendFailed,
		password.ErrValidateFailed:

		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", jsonData)
}
