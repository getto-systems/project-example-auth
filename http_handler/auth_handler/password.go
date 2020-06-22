package auth_handler

import (
	"github.com/getto-systems/project-example-id/user/authenticate"

	"github.com/getto-systems/project-example-id/data"
)

type PasswordInput struct {
	Path     string `json:"path"`
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type PasswordParam struct {
	UserID   data.UserID
	Password data.RawPassword
	Resource data.Resource
}

type PasswordHandler struct {
	AuthHandler

	AuthenticatorFactory authenticate.PasswordAuthenticatorFactory
}

func (h PasswordHandler) Handle() {
	h.Logger.Debugf(h.Request, "handling auth/password")

	param, err := h.param()
	if err != nil {
		h.errorResponse(err)
		return
	}

	token, err := h.AuthenticatorFactory.New(param.UserID, h.Request).MatchPassword(param.Password)
	if err != nil {
		h.errorResponse(err)
		return
	}

	authorizer := h.AuthorizerFactory.New(h.Request)
	if err != nil {
		h.errorResponse(err)
		return
	}

	ticket, err := authorizer.HasEnoughPermission(token, param.Resource)
	if err != nil {
		h.errorResponse(err)
		return
	}

	h.response(ticket, token)
}

func (h PasswordHandler) param() (PasswordParam, error) {
	var input PasswordInput
	err := h.parseBody(&input)
	if err != nil {
		return PasswordParam{}, err
	}

	return PasswordParam{
		UserID:   data.UserID(input.UserID),
		Password: data.RawPassword(input.Password),
		Resource: data.Resource{
			Path: data.Path(input.Path),
		},
	}, nil
}
