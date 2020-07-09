package auth_handler

import (
	"github.com/getto-systems/project-example-id/authenticate"

	"github.com/getto-systems/project-example-id/data"
)

type PasswordInput struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type PasswordParam struct {
	UserID   data.UserID
	Password data.RawPassword
}

type AuthByPassword struct {
	AuthHandler

	Auth authenticate.AuthByPassword
}

func (h AuthByPassword) Handle() {
	h.Logger.DebugMessage(&h.Request, "handling auth/by_password")

	param, err := h.param()
	if err != nil {
		h.errorResponse(err)
		return
	}

	ticket, signedTicket, err := h.Auth.Authenticate(h.Request, param.UserID, param.Password)
	if err != nil {
		h.errorResponse(err)
		return
	}

	h.response(ticket, signedTicket)
}

func (h AuthByPassword) param() (PasswordParam, error) {
	var input PasswordInput
	err := h.parseBody(&input)
	if err != nil {
		return PasswordParam{}, err
	}

	return PasswordParam{
		UserID:   data.UserID(input.UserID),
		Password: data.RawPassword(input.Password),
	}, nil
}
