package password_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"
)

type getLoginInput struct {
	Nonce string `json:"nonce"`
}

type registerInput struct {
	Nonce       string `json:"nonce"`
	LoginID     string `json:"login_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling Password/GetLogin")

	ticket, nonce, err := getLoginParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	login, err := h.password.GetLogin(request, ticket, nonce)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	loginResponse(w, login)
}

type loginResponseBody struct {
	LoginID string `json:"login_id"`
}

func loginResponse(w http.ResponseWriter, login password.Login) {
	http_handler.JsonResponse(w, http.StatusOK, loginResponseBody{
		LoginID: string(login.ID()),
	})
}

func getLoginParam(r *http.Request, logger http_handler.Logger) (_ ticket.Ticket, _ ticket.Nonce, err error) {
	var input getLoginInput
	err = http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return
	}

	nonce := ticket.Nonce(input.Nonce)

	ticket, err := http_handler.TicketCookie(r)
	if err != nil {
		return
	}

	return ticket, nonce, nil
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling Password/Register")

	ticket, nonce, login, password, err := registerParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	err = h.password.Register(request, ticket, nonce, login, password)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	h.response.OK(w)
}

func registerParam(r *http.Request, logger http_handler.Logger) (_ ticket.Ticket, _ ticket.Nonce, _ password.Login, _ password.RegisterParam, err error) {
	var input registerInput
	err = http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return
	}

	nonce := ticket.Nonce(input.Nonce)
	loginID := password.LoginID(input.LoginID)

	ticket, err := http_handler.TicketCookie(r)
	if err != nil {
		return
	}

	login := password.NewLogin(loginID)
	param := password.RegisterParam{
		OldPassword: password.RawPassword(input.OldPassword),
		NewPassword: password.RawPassword(input.NewPassword),
	}

	return ticket, nonce, login, param, nil
}
