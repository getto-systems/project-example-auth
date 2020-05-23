package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/user"
)

type PasswordHandler struct {
	CookieDomain  CookieDomain
	Authenticator auth.PasswordAuthenticator
}

func (h PasswordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param, err := passwordParam(r)
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	info, err := auth.Password(h.Authenticator, param, func(ticket user.Ticket, token auth.Token) {
		setAuthTokenCookie(w, h.CookieDomain, ticket, token)
	})
	if err != nil {
		w.WriteHeader(httpStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", info)
}

type PasswordInput struct {
	Path         string `json:"path"`
	UserID       string `json:"user_id"`
	UserPassword string `json:"password"`
}

func passwordParam(r *http.Request) (auth.PasswordParam, error) {
	var nullParam auth.PasswordParam

	if r.Body == nil {
		return nullParam, ErrBodyNotSent
	}

	var input PasswordInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return nullParam, ErrBodyParseFailed
	}

	return auth.PasswordParam{
		UserID:   user.UserID(input.UserID),
		Password: user.Password(input.UserPassword),
		Path:     user.Path(input.Path),
	}, nil
}
