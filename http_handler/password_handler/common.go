package password_handler

import (
	"net/http"

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

func (h Handler) errorResponse(w http.ResponseWriter, err error) {
	if h.password.Error().InputError(err) {
		h.response.BadRequest(w, err)
		return
	}

	if h.password.Error().AuthError(err) {
		h.response.ResetCookie(w)
		h.response.Unauthorized(w, err)
		return
	}

	if h.password.Error().ResetError(err) {
		h.response.Unauthorized(w, err)
		return
	}

	h.response.InternalServerError(w, err)
}
