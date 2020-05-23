package auth

import (
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

const COOKIE_AUTH_TOKEN = "Getto-Example-Auth-Token"

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

func setAuthTokenCookie(w http.ResponseWriter, cookieDomain CookieDomain, ticket user.Ticket, token auth.Token) {
	setter := CookieSetter{
		w,
		cookieDomain,
		ticket,
	}
	setter.setTicketCookie(token.TicketToken)
	setter.setAwsCloudFrontCookie(token.AwsCloudFrontToken)
}

func getTicketCookie(r *http.Request) (token.TicketToken, error) {
	cookie, err := r.Cookie(COOKIE_AUTH_TOKEN)
	if err != nil {
		return nil, err
	}

	return token.TicketToken(cookie.Value), nil
}

func (setter CookieSetter) setTicketCookie(ticketToken token.TicketToken) {
	setter.setCookie(&Cookie{
		Name:  COOKIE_AUTH_TOKEN,
		Value: string(ticketToken),
	})
}

func (setter CookieSetter) setAwsCloudFrontCookie(token token.AwsCloudFrontToken) {
	setter.setCookie(&Cookie{
		Name:  "CloudFront-Key-Pair-Id",
		Value: string(token.KeyPairID),
	})

	setter.setCookie(&Cookie{
		Name:  "CloudFront-Policy",
		Value: string(token.Policy),
	})

	setter.setCookie(&Cookie{
		Name:  "CloudFront-Signature",
		Value: string(token.Signature),
	})
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
