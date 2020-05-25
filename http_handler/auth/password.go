package auth

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/user"
)

type PasswordInput struct {
	Path     string `json:"path"`
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type PasswordResponse struct {
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
	AppToken string   `json:"app_token"`
}

type PasswordHandler struct {
	CookieDomain  CookieDomain
	Authenticator auth.PasswordAuthenticator
}

func (h PasswordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	logger := h.Authenticator.Logger()

	logger.Debug("auth password handling...")

	param, err := h.passwordParam(r)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("body parsed: %v", param)

	appToken, err := auth.Password(h.Authenticator, param, func(ticket user.Ticket, token auth.Token) {
		logger.Debugf("set ticket cookie: %v; %v", ticket, token)
		setAuthTokenCookie(w, h.CookieDomain, ticket, token)
	})
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Auditf("auth password success: ", param)

	jsonResponse(w, PasswordResponse{
		UserID:   string(appToken.UserID),
		Roles:    []string(appToken.Roles),
		AppToken: appToken.Token,
	})
}

func (h PasswordHandler) passwordParam(r *http.Request) (auth.PasswordParam, error) {
	logger := h.Authenticator.Logger()

	var nullParam auth.PasswordParam

	if r.Body == nil {
		logger.Info("body not sent error")
		return nullParam, ErrBodyNotSent
	}

	var input PasswordInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Info("body parse error")
		return nullParam, ErrBodyParseFailed
	}

	return auth.PasswordParam{
		UserID:   user.UserID(input.UserID),
		Password: user.Password(input.Password),
		Path:     user.Path(input.Path),
	}, nil
}
