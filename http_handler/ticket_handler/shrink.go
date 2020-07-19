package ticket_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/ticket"
)

type shrinkInput struct {
	Nonce string `json:"nonce"`
}

func (h Handler) Shrink(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling ticket/shrink")

	ticket, nonce, err := shrinkParam(r, logger)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	err = h.ticket.Shrink(request, ticket, nonce)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	h.response.ResetCookie(w)
	h.response.OK(w)
}

func shrinkParam(r *http.Request, logger http_handler.Logger) (ticket.Ticket, ticket.Nonce, error) {
	var input shrinkInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return nil, "", err
	}

	nonce := ticket.Nonce(input.Nonce)

	ticket, err := http_handler.TicketCookie(r)
	if err != nil {
		return nil, "", err
	}

	return ticket, nonce, nil
}
