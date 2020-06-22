package auth_handler

import (
	"github.com/getto-systems/project-example-id/user/authenticate"

	"github.com/getto-systems/project-example-id/data"
)

type RenewInput struct {
	Path string `json:"path"`
}

type RenewParam struct {
	Resource data.Resource
	Token    data.Token
}

type RenewHandler struct {
	AuthHandler

	AuthenticatorFactory authenticate.RenewAuthenticatorFactory
}

func (h RenewHandler) Handle() {
	h.Logger.Debugf(h.Request, "handling auth/renew")

	param, err := h.param()
	if err != nil {
		h.errorResponse(err)
		return
	}

	authorizer := h.AuthorizerFactory.New(h.Request)

	ticket, err := authorizer.HasEnoughPermission(param.Token, param.Resource)
	if err != nil {
		h.errorResponse(err)
		return
	}

	token, err := h.AuthenticatorFactory.New(ticket, h.Request).RenewTicket(ticket)
	if err != nil {
		h.errorResponse(err)
		return
	}

	h.response(ticket, token)
}

func (h RenewHandler) param() (RenewParam, error) {
	var input RenewInput
	err := h.parseBody(&input)
	if err != nil {
		return RenewParam{}, err
	}

	token, err := h.getToken()
	if err != nil {
		h.Logger.Debugf(h.Request, "token cookie not found error: %s", err)
		return RenewParam{}, ErrTokenCookieNotFound
	}

	return RenewParam{
		Token: token,
		Resource: data.Resource{
			Path: data.Path(input.Path),
		},
	}, nil
}
