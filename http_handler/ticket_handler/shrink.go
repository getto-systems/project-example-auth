package ticket_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/ticket"
)

type ShrinkInput struct {
	Nonce string `json:"nonce"`
}

type Shrink struct {
	logger   http_handler.RequestLogger
	response http_handler.Response
	shrinker ticket.TicketShrinker
}

func NewShrink(
	logger http_handler.RequestLogger,
	response http_handler.Response,
	shrinker ticket.TicketShrinker,
) Shrink {
	return Shrink{
		logger:   logger,
		response: response,
		shrinker: shrinker,
	}
}

func (h Shrink) Handle(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling ticket/shrink")

	ticket, nonce, err := h.param(r, logger)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	err = h.shrinker.Shrink(request, ticket, nonce)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	h.response.ResetCookie(w)
	h.response.OK(w)
}

func (h Shrink) param(r *http.Request, logger http_handler.Logger) (ticket.Ticket, ticket.Nonce, error) {
	var input ShrinkInput
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
