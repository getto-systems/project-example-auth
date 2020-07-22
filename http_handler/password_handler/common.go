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
	switch err {
	case password.ErrPasswordIsEmpty,
		password.ErrPasswordIsTooLong:

		h.response.BadRequest(w, err)

	case password.ErrPasswordNotFound,
		password.ErrLoginNotFound:

		h.response.ResetCookie(w)
		h.response.Unauthorized(w, err)

	case password.ErrResetTokenNotFound,
		password.ErrResetTokenNotFound,
		password.ErrResetTokenUserNotMatched,
		password.ErrResetTokenAlreadyExpired:

		h.response.Unauthorized(w, err)

	default:
		h.response.InternalServerError(w, err)
	}
}
