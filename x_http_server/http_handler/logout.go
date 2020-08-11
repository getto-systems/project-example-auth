package http_handler

import (
	"github.com/getto-systems/project-example-id"

	"github.com/getto-systems/project-example-id/request"
)

type Logout struct {
	Handler
}

func NewLogout(handler Handler) Logout {
	return Logout{Handler: handler}
}

func (handler Logout) handler() _usecase.LogoutHandler {
	return handler
}

func (handler Logout) LogoutRequest() (request.Request, error) {
	return handler.newRequest(), nil
}

func (handler Logout) LogoutResponse(err error) {
	if err != nil {
		handler.errorResponse(err)
		return
	}

	handler.ok("OK")
}