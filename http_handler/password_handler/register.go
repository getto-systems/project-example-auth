package password_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"
)

type registerInput struct {
	Nonce       string `json:"nonce"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling password/register")

	ticket, nonce, password, err := registerParam(r, logger)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	err = h.password.Register(request, ticket, nonce, password)
	if err != nil {
		h.response.ResetCookie(w)
		h.response.Error(w, err)
		return
	}

	h.response.OK(w)
}

func registerParam(r *http.Request, logger http_handler.Logger) (
	ticket.Ticket,
	ticket.Nonce,
	password.RegisterParam,
	error,
) {
	var input registerInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return nil, "", password.RegisterParam{}, err
	}

	nonce := ticket.Nonce(input.Nonce)

	ticket, err := http_handler.TicketCookie(r)
	if err != nil {
		return nil, "", password.RegisterParam{}, err
	}

	return ticket, nonce, password.RegisterParam{
		OldPassword: password.RawPassword(input.OldPassword),
		NewPassword: password.RawPassword(input.NewPassword),
	}, nil
}
