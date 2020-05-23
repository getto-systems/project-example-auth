package auth

import (
	"net/http"

	"github.com/getto-systems/project-example-id/token"
)

const COOKIE_AUTH_TOKEN = "Getto-Example-Auth-Token"

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
