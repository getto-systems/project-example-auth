package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/user"
)

type RenewHandler struct {
	CookieDomain  CookieDomain
	Authenticator auth.RenewAuthenticator
}

func (h RenewHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	logger := h.Authenticator.Logger()

	logger.Debug("auth renew handling...")

	param, err := h.renewParam(r)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("body parsed: %v", param)

	info, err := auth.Renew(h.Authenticator, param, func(ticket user.Ticket, token auth.Token) {
		logger.Debugf("set ticket cookie: %v; %v", ticket, token)
		setAuthTokenCookie(w, h.CookieDomain, ticket, token)
	})
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	logger.Debugf("auth renew ok: %v", param)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", info)
}

type RenewInput struct {
	Path string `json:"path"`
}

func (h RenewHandler) renewParam(r *http.Request) (auth.RenewParam, error) {
	logger := h.Authenticator.Logger()

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

	ticketToken, err := getTicketCookie(r)
	if err != nil {
		logger.Info("ticket cookie not sent error")
		return nullParam, ErrTicketCookieNotSent
	}

	return auth.RenewParam{
		TicketToken: ticketToken,
		Path:        user.Path(input.Path),
	}, nil
}
