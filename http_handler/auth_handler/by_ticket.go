package auth_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/authenticate"

	"github.com/getto-systems/project-example-id/data"
)

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

	signedTicket, err := h.param(r, logger)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	ticket, signedTicket, err := h.auth.Authenticate(request, signedTicket)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	h.response.write(w, ticket, signedTicket, logger)
}

func (h AuthByTicketHandler) param(r *http.Request, logger http_handler.Logger) (data.SignedTicket, error) {
	signedTicket, err := SignedTicket(r)
	if err != nil {
		logger.DebugError("signed ticket cookie not found error: %s", err)
		return nil, ErrSignedTicketCookieNotFound
	}

	return signedTicket, nil
}
