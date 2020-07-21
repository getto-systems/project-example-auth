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

	logger.DebugMessage("handling password/validate")

	login, password, err := validateParam(r, logger)
	if err != nil {
		h.response.Error(w, err)
		return
	}

	ticket, nonce, apiToken, contentToken, expires, err := h.password.Validate(request, login, password)
	if err != nil {
		h.response.ResetCookie(w)
		h.response.Error(w, err)
		return
	}

	h.response.Authenticated(w, ticket, nonce, apiToken, contentToken, expires, logger)
}

func validateParam(r *http.Request, logger http_handler.Logger) (password.Login, password.RawPassword, error) {
	var input validateInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return password.Login{}, password.RawPassword(""), err
	}

	login := password.NewLogin(password.LoginID(input.LoginID))

	return login, password.RawPassword(input.Password), nil
}
