package password_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

type RegisterInput struct {
	Nonce       string `json:"nonce"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type Register struct {
	logger   http_handler.RequestLogger
	response http_handler.Response
	register password.PasswordRegister
}

func NewRegister(
	logger http_handler.RequestLogger,
	response http_handler.Response,
	register password.PasswordRegister,
) Register {
	return Register{
		logger:   logger,
		response: response,
		register: register,
	}
}

func (h Register) Handle(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling password/register")

	ticket, nonce, password, err := h.param(r, logger)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	err = h.register.Register(request, ticket, nonce, password)
	if err != nil {
		h.response.ResetCookie(w)
		h.response.Error(w, err)
		return
	}

	h.response.OK(w)
}

func (h Register) param(r *http.Request, logger http_handler.Logger) (ticket.Ticket, ticket.Nonce, password.PasswordRegisterParam, error) {
	var input RegisterInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return nil, "", password.PasswordRegisterParam{}, err
	}

	nonce := ticket.Nonce(input.Nonce)

	ticket, err := http_handler.TicketCookie(r)
	if err != nil {
		return nil, "", password.PasswordRegisterParam{}, err
	}

	return ticket, nonce, password.PasswordRegisterParam{
		OldPassword: data.RawPassword(input.OldPassword),
		NewPassword: data.RawPassword(input.NewPassword),
	}, nil
}
