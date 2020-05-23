package auth

import (
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

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

func (setter CookieSetter) setCookie(cookie *Cookie) {
	http.SetCookie(setter.ResponseWriter, &http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,

		Domain:  string(setter.CookieDomain),
		Path:    "/",
		Expires: setter.Ticket.Expires(),

		Secure:   true,
		HttpOnly: true,
	})
}

func (setter CookieSetter) setAuthTokenCookie(token auth.Token) {
	setter.setTicketCookie(token.TicketToken)
	setter.setAwsCloudFrontCookie(token.AwsCloudFrontToken)
}
