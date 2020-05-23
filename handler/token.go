package handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/auth"
)

const COOKIE_AUTH_TOKEN = "Getto-Example-Auth-Token"

func ParseCookieToken(h AuthHandler, r *http.Request, path auth.Path) (*auth.Ticket, error) {
	cookie, err := r.Cookie(COOKIE_AUTH_TOKEN)
	if err != nil {
		return nil, err
	}

	return h.Tokener().Parse(auth.TicketToken(cookie.Value), path)
}

func SetTicketCookie(h AuthHandler, w http.ResponseWriter, ticket *auth.Ticket) error {
	value, err := h.Tokener().TicketToken(ticket)
	if err != nil {
		return err
	}

	SetCookie(h, w, ticket, &Cookie{
		Name:  COOKIE_AUTH_TOKEN,
		Value: string(value),
	})

	return nil
}
