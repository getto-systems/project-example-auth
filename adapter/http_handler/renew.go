package http_handler

import (
	"github.com/getto-systems/project-example-id/client"

	"github.com/getto-systems/project-example-id/request"
)

type Renew struct {
	Handler
}

func NewRenew(handler Handler) Renew {
	return Renew{Handler: handler}
}

func (handler Renew) handler() client.RenewHandler {
	return handler
}

func (handler Renew) RenewRequest() (request.Request, error) {
	return handler.newRequest(), nil
}

func (handler Renew) RenewResponse(err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	handler.ok("OK")
}
