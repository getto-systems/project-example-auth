package _usecase

import (
	"github.com/getto-systems/project-example-id/request"
)

type (
	renewTestHandler struct {
		*commonTestHandler
	}
)

func newRenewHandler(handler *commonTestHandler) renewTestHandler {
	return renewTestHandler{
		commonTestHandler: handler,
	}
}

func (handler renewTestHandler) handler() RenewHandler {
	return handler
}
func (handler renewTestHandler) RenewRequest() (request.Request, error) {
	return handler.newRequest(), nil
}
func (handler renewTestHandler) RenewResponse(err error) {
	handler.setError(err)
}
