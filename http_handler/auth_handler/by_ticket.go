package auth_handler

import (
	"github.com/getto-systems/project-example-id/authenticate"

	"github.com/getto-systems/project-example-id/data"
)

type TicketParam struct {
	SignedTicket data.SignedTicket
}

type AuthByTicket struct {
	AuthHandler

	Auth authenticate.AuthByTicket
}

func (h AuthByTicket) Handle() {
	h.Logger.DebugMessage(&h.Request, "handling auth/by_ticket")

	param, err := h.param()
	if err != nil {
		h.errorResponse(err)
		return
	}

	ticket, signedTicket, err := h.Auth.Authenticate(h.Request, param.SignedTicket)
	if err != nil {
		h.errorResponse(err)
		return
	}

	h.response(ticket, signedTicket)
}

func (h AuthByTicket) param() (TicketParam, error) {
	signedTicket, err := h.getSignedTicket()
	if err != nil {
		h.Logger.DebugError(&h.Request, "signed ticket cookie not found error: %s", err)
		return TicketParam{}, ErrSignedTicketCookieNotFound
	}

	return TicketParam{
		SignedTicket: signedTicket,
	}, nil
}
