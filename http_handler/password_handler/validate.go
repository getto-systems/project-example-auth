package password_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"
)

type validateInput struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

func (h Handler) Validate(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling Password/Validate")

	login, raw, err := validateParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	ticket, nonce, apiToken, contentToken, expires, err := h.password.Validate(request, login, raw)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	h.response.Authenticated(w, ticket, nonce, apiToken, contentToken, expires, logger)
}

func validateParam(r *http.Request, logger http_handler.Logger) (_ password.Login, _ password.RawPassword, err error) {
	var input validateInput
	err = http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return
	}

	login := password.NewLogin(password.LoginID(input.LoginID))
	raw := password.RawPassword(input.Password)

	return login, raw, nil
}
