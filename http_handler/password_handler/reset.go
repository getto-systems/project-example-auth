package password_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/password"
)

type createResetSessionInput struct {
	LoginID string `json:"login_id"`
}

type getResetStatusInput struct {
	ResetSessionID string `json:"reset_id"`
}

type resetInput struct {
	ResetToken string `json:"reset_token"`
	LoginID    string `json:"login_id"`
	Password   string `json:"password"`
}

type resetResponseBody struct {
	ResetSessionID string `json:"reset_id"`
}

type resetStatusResponseBody struct {
	RequestedAt    string `json:"requested_at"`
	DeliveringAt   string `json:"delivering_at"`
	DeliverdAt     string `json:"delivered_at"`
	Destination    string `json:"destination"`
	DeliverErrorAt string `json:"deliver_error_at"`
	DeliverError   string `json:"deliver_error"`
}

func (h Handler) CreateResetSession(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling Password/CreateResetSession")

	login, err := createResetSessionParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	reset, err := h.password.CreateResetSession(request, login)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	resetResponse(w, reset)
}

func createResetSessionParam(r *http.Request, logger http_handler.Logger) (_ password.Login, err error) {
	var input createResetSessionInput
	err = http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return
	}

	login := password.NewLogin(password.LoginID(input.LoginID))

	return login, nil
}

func resetResponse(w http.ResponseWriter, session password.ResetSession) {
	http_handler.JsonResponse(w, http.StatusOK, resetResponseBody{
		ResetSessionID: string(session.ID()),
	})
}

func (h Handler) GetResetStatus(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling Password/GetResetStatus")

	session, err := getResetStatusParam(r, logger)
	if err != nil {
		h.response.BadRequest(w, err)
		return
	}

	status, err := h.password.GetResetStatus(request, session)
	if err != nil {
		h.errorResponse(w, err)
		return
	}

	resetStatusResponse(w, status)
}

func getResetStatusParam(r *http.Request, logger http_handler.Logger) (_ password.ResetSession, err error) {
	var input getResetStatusInput
	err = http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return
	}

	return password.NewResetSession(password.ResetSessionID(input.ResetSessionID)), nil
}

func resetStatusResponse(w http.ResponseWriter, reset password.ResetStatus) {
	// TODO ステータスをちゃんとする
	http_handler.JsonResponse(w, http.StatusOK, resetStatusResponseBody{})
}

func (h Handler) Reset(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling Password/Reset")

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

func resetParam(r *http.Request, logger http_handler.Logger) (_ password.Login, _ password.ResetToken, _ password.RawPassword, err error) {
	var input resetInput
	err = http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return
	}

	login := password.NewLogin(password.LoginID(input.LoginID))
	token := password.ResetToken(input.ResetToken)
	raw := password.RawPassword(input.Password)

	return login, token, raw, nil
}
