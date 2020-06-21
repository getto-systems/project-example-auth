package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/user/authenticate"
	"github.com/getto-systems/project-example-id/user/authorize"

	"github.com/getto-systems/project-example-id/data"

	"errors"
	"fmt"
)

const COOKIE_AUTH_TOKEN = "Getto-Example-Auth-Token"

type AuthHandler struct {
	Logger http_handler.Logger

	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request

	CookieDomain CookieDomain

	AwsCloudFrontIssuer AwsCloudFrontIssuer
	AppIssuer           AppIssuer

	Request data.Request

	AuthorizerFactory authorize.AuthorizerFactory
}

type CookieDomain string

var (
	ErrEmptyBody           = errors.New("empty body")
	ErrBodyParseFailed     = errors.New("body parse failed")
	ErrTokenCookieNotFound = errors.New("token cookie not found")
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
		h.Logger.Debugf(h.Request, "empty body error")
		return ErrEmptyBody
	}

	err := json.NewDecoder(h.HttpRequest.Body).Decode(&input)
	if err != nil {
		h.Logger.Debugf(h.Request, "body parse error: %s", err)
		return ErrBodyParseFailed
	}

	return nil
}

func (h AuthHandler) response(ticket data.Ticket, token data.Token) {
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

	h.setTokenCookie(token, ticket.Expires)
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
		ErrTokenCookieNotFound,
		authenticate.ErrPasswordMatchFailed,
		authenticate.ErrTicketRenewFailed,
		authorize.ErrAuthorizeTokenParseFailed:

		h.HttpResponseWriter.WriteHeader(http.StatusUnauthorized)

	case authorize.ErrAuthorizeFailed:
		h.HttpResponseWriter.WriteHeader(http.StatusForbidden)

	default:
		h.Logger.Debugf(h.Request, "internal server error: %s", err)
		h.HttpResponseWriter.WriteHeader(http.StatusInternalServerError)
	}
}

func (h AuthHandler) getToken() (data.Token, error) {
	cookie, err := h.HttpRequest.Cookie(COOKIE_AUTH_TOKEN)
	if err != nil {
		return nil, err
	}

	return data.Token(cookie.Value), nil
}

func (h AuthHandler) setTokenCookie(token data.Token, expires data.Expires) {
	http.SetCookie(h.HttpResponseWriter, &http.Cookie{
		Name:  COOKIE_AUTH_TOKEN,
		Value: string(token),

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
