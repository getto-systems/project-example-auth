package http_handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/getto-systems/project-example-id"

	"github.com/getto-systems/project-example-id/credential"
)

var (
	errTicketTokenNotFound = errors.New("ticket token not found")
)

const (
	COOKIE_TICKET    = "GETTO-EXAMPLE-ID-TicketToken"
	COOKIE_API_TOKEN = "GETTO-EXAMPLE-ID-ApiToken"

	HEADER_NONCE     = "X-GETTO-EXAMPLE-ID-TicketNonce"
	HEADER_API_ROLES = "X-GETTO-EXAMPLE-ID-ApiRoles"
)

type (
	Cookie struct {
		domain CookieDomain
	}

	CookieDomain string

	CredentialHandler struct {
		cookie Cookie

		httpResponseWriter http.ResponseWriter
		httpRequest        *http.Request
	}
)

func NewCredentialHandler(domain CookieDomain, w http.ResponseWriter, r *http.Request) CredentialHandler {
	return CredentialHandler{
		cookie: newCookie(domain),

		httpResponseWriter: w,
		httpRequest:        r,
	}
}

func newCookie(domain CookieDomain) Cookie {
	return Cookie{
		domain: domain,
	}
}

func (handler CredentialHandler) handler() _usecase.CredentialHandler {
	return handler
}

func (handler CredentialHandler) GetTicketNonceAndSignature() (_ credential.TicketNonce, _ credential.TicketSignature, err error) {
	cookie, err := handler.httpRequest.Cookie(COOKIE_TICKET)
	if err != nil {
		err = errTicketTokenNotFound
		return
	}

	nonce := handler.httpRequest.Header.Get(HEADER_NONCE)

	return credential.TicketNonce(nonce), credential.TicketSignature(cookie.Value), nil
}

func (handler CredentialHandler) SetCredential(credential credential.Credential) {
	handler.setTicket(credential.TicketToken())
	handler.setApiToken(credential.ApiToken())
	handler.setContentToken(credential.ContentToken())
}

func (handler CredentialHandler) ClearCredential() {
	handler.clearTicket()
	handler.clearContentToken()
}

func (handler CredentialHandler) setTicket(ticket credential.TicketToken) {
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    COOKIE_TICKET,
		Value:   string(ticket.Signature()),
		Expires: time.Time(ticket.Expires()),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	handler.httpResponseWriter.Header().Set(
		HEADER_NONCE,
		string(ticket.Nonce()),
	)
}

func (handler CredentialHandler) clearTicket() {
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:   COOKIE_TICKET,
		Value:  "",
		MaxAge: -1,

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (handler CredentialHandler) setApiToken(apiToken credential.ApiToken) {
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    COOKIE_API_TOKEN,
		Value:   string(apiToken.Signature()),
		Expires: time.Time(apiToken.Expires()),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	jsonData, err := json.Marshal(apiToken.ApiRoles())
	if err != nil {
		return
	}

	handler.httpResponseWriter.Header().Set(
		HEADER_API_ROLES,
		base64.StdEncoding.EncodeToString(jsonData),
	)
}

func (handler CredentialHandler) setContentToken(contentToken credential.ContentToken) {
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    "CloudFront-Key-Pair-Id",
		Value:   string(contentToken.KeyID()),
		Expires: time.Time(contentToken.Expires()),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    "CloudFront-Policy",
		Value:   string(contentToken.Policy()),
		Expires: time.Time(contentToken.Expires()),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    "CloudFront-Signature",
		Value:   string(contentToken.Signature()),
		Expires: time.Time(contentToken.Expires()),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (handler CredentialHandler) clearContentToken() {
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:   "CloudFront-Key-Pair-Id",
		Value:  "",
		MaxAge: -1,

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:   "CloudFront-Policy",
		Value:  "",
		MaxAge: -1,

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:   "CloudFront-Signature",
		Value:  "",
		MaxAge: -1,

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
