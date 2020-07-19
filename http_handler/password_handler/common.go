package password_handler

import (
	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"
)

type Handler struct {
	logger   http_handler.RequestLogger
	response http_handler.Response

	password password.Usecase
}

func NewHandler(
	logger http_handler.RequestLogger,
	response http_handler.Response,
	password password.Usecase,
) Handler {
	return Handler{
		logger:   logger,
		response: response,

		password: password,
	}
}
