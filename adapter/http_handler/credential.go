package http_handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	gotime "time"

	"github.com/getto-systems/project-example-id/client"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/time"
)

var (
	errTicketTokenNotFound = errors.New("ticket token not found")
)

const (
	COOKIE_TICKET = "GETTO-EXAMPLE-ID.TicketToken"

	HEADER_NONCE     = "X-GETTO-EXAMPLE-ID.TicketNonce"
	HEADER_API_TOKEN = "X-GETTO-EXAMPLE-ID.TicketNonce"
	HEADER_API_ROLES = "X-GETTO-EXAMPLE-ID.TicketNonce"
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

func (handler CredentialHandler) handler() client.CredentialHandler {
	return handler
}

func (handler CredentialHandler) GetTicket() (_ api_token.Ticket, err error) {
	cookie, err := handler.httpRequest.Cookie(COOKIE_TICKET)
	if err != nil {
		err = errTicketTokenNotFound
		return
	}

	nonce := handler.httpRequest.Header.Get(HEADER_NONCE)

	return api_token.NewTicket(api_token.TicketSignature(cookie.Value), api_token.TicketNonce(nonce)), nil
}

func (handler CredentialHandler) SetCredential(credential data.Credential) {
	handler.setTicket(credential.Ticket(), credential.Expires())
	handler.setApiToken(credential.ApiToken())
	handler.setContentToken(credential.ContentToken(), credential.Expires())
}

func (handler CredentialHandler) ClearCredential() {
	handler.clearTicket()
	handler.clearContentToken()
}

func (handler CredentialHandler) setTicket(ticket api_token.Ticket, expires time.Expires) {
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    COOKIE_TICKET,
		Value:   string(ticket.Signature()),
		Expires: gotime.Time(expires),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
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

func (handler CredentialHandler) setApiToken(apiToken api_token.ApiToken) {
	handler.httpResponseWriter.Header().Set(
		HEADER_API_TOKEN,
		string(apiToken.Signature()),
	)

	jsonData, err := json.Marshal(apiToken.ApiRoles())
	if err != nil {
		handler.httpResponseWriter.Header().Set(
			HEADER_API_ROLES,
			base64.StdEncoding.EncodeToString(jsonData),
		)
	}
}

func (handler CredentialHandler) setContentToken(contentToken api_token.ContentToken, expires time.Expires) {
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    "CloudFront-Key-Pair-Id",
		Value:   string(contentToken.KeyID()),
		Expires: gotime.Time(expires),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    "CloudFront-Policy",
		Value:   string(contentToken.Policy()),
		Expires: gotime.Time(expires),

		Domain: string(handler.cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(handler.httpResponseWriter, &http.Cookie{
		Name:    "CloudFront-Signature",
		Value:   string(contentToken.Signature()),
		Expires: gotime.Time(expires),

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
