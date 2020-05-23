package renew

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/getto-systems/project-example-id/http_handler"
	auth_handler "github.com/getto-systems/project-example-id/http_handler/auth"

	"github.com/getto-systems/project-example-id/auth"
	auth_renew "github.com/getto-systems/project-example-id/auth/renew"

	"github.com/getto-systems/project-example-id/user"
)

type Handler struct {
	CookieDomain  http_handler.CookieDomain
	Authenticator auth_renew.Authenticator
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param, err := param(r)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	info, err := auth_renew.Renew(h.Authenticator, param, func(ticket user.Ticket, token auth.Token) {
		setter := http_handler.CookieSetter{
			ResponseWriter: w,
			CookieDomain:   h.CookieDomain,
			Ticket:         ticket,
		}
		auth_handler.SetAuthTokenCookie(setter, token)
	})
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", info)
}

type Input struct {
	Path string `json:"path"`
}

var ErrBodyNotSent = errors.New("body not sent")
var ErrBodyParseFailed = errors.New("body parse failed")
var ErrTicketCookieNotSent = errors.New("ticket cookie not sent")

func param(r *http.Request) (auth_renew.RenewParam, error) {
	var nullParam auth_renew.RenewParam

	if r.Body == nil {
		return nullParam, ErrBodyNotSent
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return nullParam, ErrBodyParseFailed
	}

	ticketToken, err := auth_handler.GetTicketCookie(r)
	if err != nil {
		return nullParam, ErrTicketCookieNotSent
	}

	return auth_renew.RenewParam{
		TicketToken: ticketToken,
		Path:        user.Path(input.Path),
		Now:         time.Now().UTC(),
	}, nil
}

func httpStatusCode(err error) int {
	switch err {

	case ErrBodyNotSent, ErrBodyParseFailed:
		return http.StatusBadRequest

	case ErrTicketCookieNotSent, auth.ErrTicketTokenParseFailed:
		return http.StatusUnauthorized

	case auth.ErrUserAccessDenied:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
