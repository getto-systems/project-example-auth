package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/user"
)

type RenewHandler struct {
	CookieDomain  http_handler.CookieDomain
	Authenticator auth.RenewAuthenticator
}

func (h RenewHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param, err := renewParam(r)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	info, err := auth.Renew(h.Authenticator, param, func(ticket user.Ticket, token auth.Token) {
		setter := http_handler.CookieSetter{
			ResponseWriter: w,
			CookieDomain:   h.CookieDomain,
			Ticket:         ticket,
		}
		SetAuthTokenCookie(setter, token)
	})
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", info)
}

type RenewInput struct {
	Path string `json:"path"`
}

func renewParam(r *http.Request) (auth.RenewParam, error) {
	var nullParam auth.RenewParam

	if r.Body == nil {
		return nullParam, ErrBodyNotSent
	}

	var input RenewInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return nullParam, ErrBodyParseFailed
	}

	ticketToken, err := GetTicketCookie(r)
	if err != nil {
		return nullParam, ErrTicketCookieNotSent
	}

	return auth.RenewParam{
		TicketToken: ticketToken,
		Path:        user.Path(input.Path),
	}, nil
}
