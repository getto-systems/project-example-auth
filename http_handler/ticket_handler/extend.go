package ticket_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/ticket"
)

type ExtendInput struct {
	Nonce string `json:"nonce"`
}

type Extend struct {
	logger   http_handler.RequestLogger
	response http_handler.Response
	extender ticket.TicketExtender
}

func NewExtend(
	logger http_handler.RequestLogger,
	response http_handler.Response,
	extender ticket.TicketExtender,
) Extend {
	return Extend{
		logger:   logger,
		response: response,
		extender: extender,
	}
}

func (h Extend) Handle(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling ticket/extend")

	ticket, nonce, err := h.param(r, logger)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	ticket, apiToken, contentToken, expires, err := h.extender.Extend(request, ticket, nonce)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	h.response.Authenticated(w, ticket, nonce, apiToken, contentToken, expires, logger)
}

func (h Extend) param(r *http.Request, logger http_handler.Logger) (ticket.Ticket, ticket.Nonce, error) {
	var input ExtendInput
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
