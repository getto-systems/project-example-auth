package http_handler

import (
	"net/http"
	"time"

	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

const COOKIE_TICKET = "GETTO-EXAMPLE-ID.Ticket"

type (
	Cookie struct {
		domain         CookieDomain
		contentTokenID ContentTokenID
	}

	CookieDomain   string
	ContentTokenID string
)

func NewCookie(domain CookieDomain, contentTokenID ContentTokenID) Cookie {
	return Cookie{
		domain:         domain,
		contentTokenID: contentTokenID,
	}
}

func TicketCookie(r *http.Request) (ticket.Ticket, error) {
	cookie, err := r.Cookie(COOKIE_TICKET)
	if err != nil {
		return nil, ErrTicketCookieNotFound
	}

	return ticket.Ticket(cookie.Value), nil
}

func (cookie Cookie) setTicket(w http.ResponseWriter, ticket ticket.Ticket, expires data.Expires) {
	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_TICKET,
		Value:   string(ticket),
		Expires: time.Time(expires),

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (cookie Cookie) resetTicket(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   COOKIE_TICKET,
		Value:  "",
		MaxAge: -1,

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (cookie Cookie) setContentToken(w http.ResponseWriter, contentToken ticket.ContentToken, expires data.Expires) {
	http.SetCookie(w, &http.Cookie{
		Name:    "CloudFront-Key-Pair-Id",
		Value:   string(cookie.contentTokenID),
		Expires: time.Time(expires),

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "CloudFront-Policy",
		Value:   string(contentToken.Policy()),
		Expires: time.Time(expires),

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "CloudFront-Signature",
		Value:   string(contentToken.Signature()),
		Expires: time.Time(expires),

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (cookie Cookie) resetContentToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "CloudFront-Key-Pair-Id",
		Value:  "",
		MaxAge: -1,

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "CloudFront-Policy",
		Value:  "",
		MaxAge: -1,

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "CloudFront-Signature",
		Value:  "",
		MaxAge: -1,

		Domain: string(cookie.domain),
		Path:   "/",

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
