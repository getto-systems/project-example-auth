package password_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"

	"github.com/getto-systems/project-example-id/data"
)

type ValidateInput struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type Validate struct {
	logger    http_handler.RequestLogger
	response  http_handler.Response
	validater password.PasswordValidater
}

func NewValidate(
	logger http_handler.RequestLogger,
	response http_handler.Response,
	validater password.PasswordValidater,
) Validate {
	return Validate{
		logger:    logger,
		response:  response,
		validater: validater,
	}
}

func (h Validate) Handle(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling password/validate")

	user, password, err := h.param(r, logger)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	ticket, nonce, apiToken, contentToken, expires, err := h.validater.Validate(request, user, password)
	if err != nil {
		h.response.ResetCookie(w)
		h.response.Error(w, err)
		return
	}

	h.response.Authenticated(w, ticket, nonce, apiToken, contentToken, expires, logger)
}

func (h Validate) param(r *http.Request, logger http_handler.Logger) (data.User, data.RawPassword, error) {
	var input ValidateInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return data.User{}, data.RawPassword(""), err
	}

	user := data.NewUser(data.UserID(input.UserID))

	return user, data.RawPassword(input.Password), nil
}
