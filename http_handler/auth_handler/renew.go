package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/applog"

	"github.com/getto-systems/project-example-id/basic"
)

type RenewHandler struct {
	CookieDomain         CookieDomain
	AuthenticatorFactory func(*http.Request) auth.RenewAuthenticator
}

type RenewInput struct {
	Path string `json:"path"`
}

type RenewResponse struct {
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
	AppToken string   `json:"app_token"`
}

func (h RenewHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authenticator := h.AuthenticatorFactory(r)

	logger := authenticator.Logger()

	logger.Debug("handling auth/renew...")

	param, err := renewParam(r, logger)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("body parsed: %v", param)

	appToken, err := auth.Renew(authenticator, param, func(token auth.Token) {
		logger.Debugf("set ticket cookie: %v", token)
		setAuthTokenCookie(w, h.CookieDomain, token)
	})
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("auth renew ok: %v", param)

	jsonResponse(w, RenewResponse{
		UserID:   string(appToken.UserID),
		Roles:    []string(appToken.Roles),
		AppToken: appToken.Token,
	})
}

func renewParam(r *http.Request, logger applog.Logger) (auth.RenewParam, error) {
	if r.Body == nil {
		logger.Info("body not sent error")
		return auth.RenewParam{}, ErrBodyNotSent
	}

	var input RenewInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Info("body parse error")
		return auth.RenewParam{}, ErrBodyParseFailed
	}

	renewToken, err := getRenewToken(r)
	if err != nil {
		logger.Info("ticket cookie not sent error")
		return auth.RenewParam{}, ErrTicketCookieNotSent
	}

	return auth.RenewParam{
		RequestedAt: http_handler.Now(),

		RenewToken: renewToken,
		Path:       basic.Path(input.Path),
	}, nil
}
