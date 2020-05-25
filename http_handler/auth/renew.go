package auth

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/logger"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/user"
)

type RenewHandler struct {
	CookieDomain         CookieDomain
	AuthenticatorFactory func(*http.Request) (auth.RenewAuthenticator, error)
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

	authenticator, err := h.AuthenticatorFactory(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger := authenticator.Logger()

	logger.Debug("auth renew handling...")

	param, err := renewParam(r, logger)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("body parsed: %v", param)

	appToken, err := auth.Renew(authenticator, param, func(ticket user.Ticket, token auth.Token) {
		logger.Debugf("set ticket cookie: %v; %v", ticket, token)
		setAuthTokenCookie(w, h.CookieDomain, ticket, token)
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

func renewParam(r *http.Request, logger logger.Logger) (auth.RenewParam, error) {
	var nullParam auth.RenewParam

	if r.Body == nil {
		logger.Info("body not sent error")
		return nullParam, ErrBodyNotSent
	}

	var input RenewInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Info("body parse error")
		return nullParam, ErrBodyParseFailed
	}

	renewToken, err := getRenewToken(r)
	if err != nil {
		logger.Info("ticket cookie not sent error")
		return nullParam, ErrTicketCookieNotSent
	}

	return auth.RenewParam{
		RenewToken: renewToken,
		Path:       user.Path(input.Path),
	}, nil
}
