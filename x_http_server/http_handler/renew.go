package http_handler

import (
	"github.com/getto-systems/project-example-id"

	"github.com/getto-systems/project-example-id/request"
)

type Renew struct {
	Handler
}

func NewRenew(handler Handler) Renew {
	return Renew{Handler: handler}
}

func (handler Renew) handler() _usecase.RenewHandler {
	return handler
}

func (handler Renew) RenewRequest() (request.Request, error) {
	return handler.newRequest(), nil
}

func (handler Renew) RenewResponse(err error) {
	if err != nil {
		switch err {
		case _usecase.ErrBadRequest:
			handler.badRequest()

		case _usecase.ErrInvalidTicket:
			handler.invalidTicket()

		default:
			handler.internalServerError()
		}
		return
	}

	handler.ok("OK")
}
