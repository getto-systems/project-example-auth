package password

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"
	auth_handler "github.com/getto-systems/project-example-id/http_handler/auth"

	"github.com/getto-systems/project-example-id/auth"
	auth_password "github.com/getto-systems/project-example-id/auth/password"

	"github.com/getto-systems/project-example-id/user"
	user_password "github.com/getto-systems/project-example-id/user/password"
)

type Handler struct {
	CookieDomain  http_handler.CookieDomain
	Authenticator auth_password.Authenticator
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param, err := param(r)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	info, err := auth_password.Auth(h.Authenticator, param, func(ticket user.Ticket, token auth.Token) {
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
	Path         string `json:"path"`
	UserID       string `json:"user_id"`
	UserPassword string `json:"password"`
}

var ErrBodyNotSent = errors.New("body not sent")
var ErrBodyParseFailed = errors.New("body parse failed")

func param(r *http.Request) (auth_password.AuthParam, error) {
	var nullParam auth_password.AuthParam

	if r.Body == nil {
		return nullParam, ErrBodyNotSent
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return nullParam, ErrBodyParseFailed
	}

	return auth_password.AuthParam{
		UserID:       user.UserID(input.UserID),
		UserPassword: user_password.UserPassword(input.UserPassword),
		Path:         user.Path(input.Path),
	}, nil
}

func httpStatusCode(err error) int {
	switch err {

	case ErrBodyNotSent, ErrBodyParseFailed:
		return http.StatusBadRequest

	case auth.ErrTicketTokenParseFailed:
		return http.StatusUnauthorized

	case auth.ErrUserAccessDenied:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
