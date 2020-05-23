package auth

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/token"
)

const COOKIE_AUTH_TOKEN = "Getto-Example-Auth-Token"

func GetTicketCookie(r *http.Request) (token.TicketToken, error) {
	cookie, err := r.Cookie(COOKIE_AUTH_TOKEN)
	if err != nil {
		return nil, err
	}

	return token.TicketToken(cookie.Value), nil
}

func SetTicketCookie(setter http_handler.CookieSetter, ticketToken token.TicketToken) {
	setter.SetCookie(&http_handler.Cookie{
		Name:  COOKIE_AUTH_TOKEN,
		Value: string(ticketToken),
	})
}
