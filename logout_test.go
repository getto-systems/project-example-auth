package auth

import (
	"github.com/getto-systems/project-example-id/request"
)

type (
	logoutTestHandler struct {
		*commonTestHandler
	}
)

func newLogoutHandler(handler *commonTestHandler) logoutTestHandler {
	return logoutTestHandler{
		commonTestHandler: handler,
	}
}

func (handler logoutTestHandler) handler() LogoutHandler {
	return handler
}
func (handler logoutTestHandler) LogoutRequest() (request.Request, error) {
	return handler.newRequest(), nil
}
func (handler logoutTestHandler) LogoutResponse(err error) {
	handler.setError(err)
}
