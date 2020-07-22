package password_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"
)

type issueResetInput struct {
	LoginID string `json:"login_id"`
}

type getResetStatusInput struct {
	ResetID string `json:"reset_id"`
}

type resetInput struct {
	ResetToken string `json:"reset_token"`
	LoginID    string `json:"login_id"`
	Password   string `json:"password"`
}

type resetResponseBody struct {
	ResetID string `json:"reset_id"`
}

type resetStatusResponseBody struct {
	RequestedAt    string `json:"requested_at"`
	DeliveringAt   string `json:"delivering_at"`
	DeliverdAt     string `json:"delivered_at"`
	Destination    string `json:"destination"`
	DeliverErrorAt string `json:"deliver_error_at"`
	DeliverError   string `json:"deliver_error"`
}

func (h Handler) IssueReset(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling password/validate")

	login, err := issueResetParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	reset, err := h.password.IssueReset(request, login)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	resetResponse(w, reset)
}

func issueResetParam(r *http.Request, logger http_handler.Logger) (password.Login, error) {
	var input issueResetInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return password.Login{}, err
	}

	login := password.NewLogin(password.LoginID(input.LoginID))

	return login, nil
}

func resetResponse(w http.ResponseWriter, reset password.Reset) {
	http_handler.JsonResponse(w, http.StatusOK, resetResponseBody{
		ResetID: string(reset.ID()),
	})
}

func (h Handler) GetResetStatus(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling password/validate")

	reset, err := getResetStatusParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	status, err := h.password.GetResetStatus(request, reset)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	resetStatusResponse(w, status)
}

func getResetStatusParam(r *http.Request, logger http_handler.Logger) (password.Reset, error) {
	var input getResetStatusInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return password.Reset{}, err
	}

	reset := password.NewReset(password.ResetID(input.ResetID))

	return reset, nil
}

func resetStatusResponse(w http.ResponseWriter, reset password.ResetStatus) {
	// TODO ステータスをちゃんとする
	http_handler.JsonResponse(w, http.StatusOK, resetStatusResponseBody{})
}

func (h Handler) Reset(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling password/validate")

	login, token, raw, err := resetParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	ticket, nonce, apiToken, contentToken, expires, err := h.password.Reset(request, login, token, raw)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	h.response.Authenticated(w, ticket, nonce, apiToken, contentToken, expires, logger)
}

func resetParam(r *http.Request, logger http_handler.Logger) (password.Login, password.ResetToken, password.RawPassword, error) {
	var input resetInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return password.Login{}, "", "", err
	}

	login := password.NewLogin(password.LoginID(input.LoginID))
	token := password.ResetToken(input.ResetToken)
	raw := password.RawPassword(input.Password)

	return login, token, raw, nil
}
