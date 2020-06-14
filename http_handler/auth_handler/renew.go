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
	AuthHandler
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

	logger := h.Logger()

	logger.Debug("handling auth/renew...")

	param, err := h.renewParam(r, logger)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("body parsed: %v", param)

	authenticator := h.AuthenticatorFactory(r)

	ticket, err := auth.Renew(authenticator, param)
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

	logger.Debugf("auth renew ok: %v", param)

	appToken, err := SerializeAppToken(h, ticket)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	jsonResponse(w, RenewResponse{
		UserID:   string(appToken.UserID),
		Roles:    []string(appToken.Roles),
		AppToken: appToken.Token,
	})
}

func (h RenewHandler) renewParam(r *http.Request, logger applog.Logger) (auth.RenewParam, error) {
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

	path := basic.Path(input.Path)

	// TODO コピペしただけ
	ticketSerializer := h.TicketSerializer()

	ticket, err := ticketSerializer.Parse(renewToken, path)
	if err != nil {
		logger.Debugf("parse token error: %s; %s / $s", err, renewToken, path)
		return auth.RenewParam{}, ErrRenewTokenParseFailed
	}

	return auth.RenewParam{
		RequestedAt: http_handler.RequestedAt(),

		Ticket: ticket,
		Path:   path,
	}, nil
}
