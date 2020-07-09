package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/authenticate"

	"github.com/getto-systems/project-example-id/data"

	"errors"
	"fmt"
)

const COOKIE_SIGNED_TICKET = "Getto-Example-Auth-Signed-Ticket"

type AuthHandler struct {
	Logger http_handler.Logger

	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request

	CookieDomain CookieDomain

	AwsCloudFrontIssuer AwsCloudFrontIssuer
	AppIssuer           AppIssuer

	Request data.Request
}

type CookieDomain string

var (
	ErrEmptyBody                  = errors.New("empty body")
	ErrBodyParseFailed            = errors.New("body parse failed")
	ErrSignedTicketCookieNotFound = errors.New("signed ticket cookie not found")
)

func Request(r *http.Request) data.Request {
	return data.Request{
		RequestedAt: http_handler.RequestedAt(),
		Route: data.Route{
			RemoteAddr: data.RemoteAddr(r.RemoteAddr),
		},
	}
}

func (h AuthHandler) parseBody(input interface{}) error {
	if h.HttpRequest.Body == nil {
		h.Logger.DebugError(&h.Request, "empty body error", nil)
		return ErrEmptyBody
	}

	err := json.NewDecoder(h.HttpRequest.Body).Decode(&input)
	if err != nil {
		h.Logger.DebugError(&h.Request, "body parse error: %s", err)
		return ErrBodyParseFailed
	}

	return nil
}

func (h AuthHandler) response(ticket data.Ticket, signedTicket data.SignedTicket) {
	awsCloudFrontTicket, err := h.AwsCloudFrontIssuer.Authorized(ticket)
	if err != nil {
		h.errorResponse(err)
		return
	}

	appToken, err := h.AppIssuer.Authorized(ticket)
	if err != nil {
		h.errorResponse(err)
		return
	}

	h.setSignedTicketCookie(signedTicket, ticket.Expires)
	h.setAwsCloudFrontCookie(awsCloudFrontTicket, ticket.Expires)

	h.jsonResponse(appToken)
}

func (h AuthHandler) jsonResponse(response interface{}) {
	h.HttpResponseWriter.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)
	if err != nil {
		h.HttpResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.HttpResponseWriter.WriteHeader(http.StatusOK)
	fmt.Fprintf(h.HttpResponseWriter, "%s", data)
}

func (h AuthHandler) errorResponse(err error) {
	h.HttpResponseWriter.Header().Set("Content-Type", "application/json")

	switch err {

	case
		ErrEmptyBody,
		ErrBodyParseFailed:

		h.HttpResponseWriter.WriteHeader(http.StatusBadRequest)

	case
		ErrSignedTicketCookieNotFound,
		authenticate.ErrTicketAuthFailed,
		authenticate.ErrPasswordAuthFailed:

		h.HttpResponseWriter.WriteHeader(http.StatusUnauthorized)

	default:
		h.Logger.DebugError(&h.Request, "internal server error: %s", err)
		h.HttpResponseWriter.WriteHeader(http.StatusInternalServerError)
	}
}

func (h AuthHandler) getSignedTicket() (data.SignedTicket, error) {
	cookie, err := h.HttpRequest.Cookie(COOKIE_SIGNED_TICKET)
	if err != nil {
		return nil, err
	}

	return data.SignedTicket(cookie.Value), nil
}

func (h AuthHandler) setSignedTicketCookie(signedTicket data.SignedTicket, expires data.Expires) {
	http.SetCookie(h.HttpResponseWriter, &http.Cookie{
		Name:  COOKIE_SIGNED_TICKET,
		Value: string(signedTicket),

		Domain:  string(h.CookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h AuthHandler) setAwsCloudFrontCookie(ticket AwsCloudFrontTicket, expires data.Expires) {
	http.SetCookie(h.HttpResponseWriter, &http.Cookie{
		Name:  "CloudFront-Key-Pair-Id",
		Value: string(ticket.KeyPairID),

		Domain:  string(h.CookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(h.HttpResponseWriter, &http.Cookie{
		Name:  "CloudFront-Policy",
		Value: string(ticket.Token.Policy),

		Domain:  string(h.CookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(h.HttpResponseWriter, &http.Cookie{
		Name:  "CloudFront-Signature",
		Value: string(ticket.Token.Signature),

		Domain:  string(h.CookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
