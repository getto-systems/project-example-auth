package http_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/user"
)

type CookieSetter struct {
	ResponseWriter http.ResponseWriter
	CookieDomain   CookieDomain
	Ticket         user.Ticket
}

type CookieDomain string

type Cookie struct {
	Name  string
	Value string
}

func (c CookieSetter) SetCookie(cookie *Cookie) {
	http.SetCookie(c.ResponseWriter, &http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,

		Domain:  string(c.CookieDomain),
		Path:    "/",
		Expires: c.Ticket.Expires(),

		Secure:   true,
		HttpOnly: true,
	})
}
