package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/applog"

	"github.com/getto-systems/project-example-id/basic"
)

type PasswordHandler struct {
	AuthHandler
	CookieDomain         CookieDomain
	AuthenticatorFactory func(*http.Request) auth.PasswordAuthenticator
}

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

func (h PasswordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	logger := h.Logger()

	logger.Debug("handling auth/password...")

	param, err := passwordParam(r, logger)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("body parsed: %v", param)

	authenticator := h.AuthenticatorFactory(r)

	ticket, err := auth.Password(authenticator, param)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	token, err := SerializeAuthToken(h, ticket)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("set ticket cookie: %v", token)
	setAuthTokenCookie(w, h.CookieDomain, token)

	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	appToken, err := SerializeAppToken(h, ticket)
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

func passwordParam(r *http.Request, logger applog.Logger) (auth.PasswordParam, error) {
	if r.Body == nil {
		logger.Info("body not sent error")
		return auth.PasswordParam{}, ErrBodyNotSent
	}

	var input PasswordInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Info("body parse error")
		return auth.PasswordParam{}, ErrBodyParseFailed
	}

	return auth.PasswordParam{
		RequestedAt: http_handler.RequestedAt(),

		UserID:   basic.UserID(input.UserID),
		Password: basic.Password(input.Password),
		Path:     basic.Path(input.Path),
	}, nil
}
