package ticket_handler

import (
	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/ticket"
)

type Handler struct {
	logger   http_handler.RequestLogger
	response http_handler.Response

	ticket ticket.Usecase
}

func NewHandler(
	logger http_handler.RequestLogger,
	response http_handler.Response,
	ticket ticket.Usecase,
) Handler {
	return Handler{
		logger:   logger,
		response: response,

		ticket: ticket,
	}
}
