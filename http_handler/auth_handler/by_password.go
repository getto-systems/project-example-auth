package auth_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

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

type AuthByPasswordHandler struct {
	logger   http_handler.RequestLogger
	response AuthResponse
	auth     authenticate.AuthByPassword
}

func NewAuthByPasswordHandler(logger http_handler.RequestLogger, response AuthResponse, auth authenticate.AuthByPassword) AuthByPasswordHandler {
	return AuthByPasswordHandler{
		logger:   logger,
		response: response,
		auth:     auth,
	}
}

func (h AuthByPasswordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	request := http_handler.Request(r)
	logger := http_handler.NewLogger(h.logger, request)

	logger.DebugMessage("handling auth/by_password")

	param, err := passwordParam(r, logger)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	ticket, signedTicket, err := h.auth.Authenticate(request, param.UserID, param.Password)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	h.response.write(w, ticket, signedTicket, logger)
}

func passwordParam(r *http.Request, logger http_handler.Logger) (PasswordParam, error) {
	var input PasswordInput
	err := http_handler.ParseBody(r, &input, logger)
	if err != nil {
		return PasswordParam{}, err
	}

	return PasswordParam{
		UserID:   data.UserID(input.UserID),
		Password: data.RawPassword(input.Password),
	}, nil
}
