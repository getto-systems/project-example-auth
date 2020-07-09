package auth_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/authenticate"

	"github.com/getto-systems/project-example-id/data"
)

type TicketParam struct {
	SignedTicket data.SignedTicket
}

type AuthByTicketHandler struct {
	logger   http_handler.RequestLogger
	response AuthResponse
	auth     authenticate.AuthByTicket
}

func NewAuthByTicketHandler(logger http_handler.RequestLogger, response AuthResponse, auth authenticate.AuthByTicket) AuthByTicketHandler {
	return AuthByTicketHandler{
		logger:   logger,
		response: response,
		auth:     auth,
	}
}

func (h AuthByTicketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling auth/by_ticket")

	param, err := ticketParam(r, logger)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	ticket, signedTicket, err := h.auth.Authenticate(request, param.SignedTicket)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	h.response.write(w, ticket, signedTicket, logger)
}

func ticketParam(r *http.Request, logger http_handler.Logger) (TicketParam, error) {
	signedTicket, err := SignedTicket(r)
	if err != nil {
		logger.DebugError("signed ticket cookie not found error: %s", err)
		return TicketParam{}, ErrSignedTicketCookieNotFound
	}

	return TicketParam{
		SignedTicket: signedTicket,
	}, nil
}
