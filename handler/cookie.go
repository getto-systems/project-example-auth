package handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/signature"
)

type Domain string

type CookieHandler interface {
	Domain() Domain
}
type CredentialHandler interface {
	CookieHandler
	Tokener() auth.Tokener
	CloudFrontSigner() signature.CloudFrontSigner
}

type Cookie struct {
	Name  string
	Value string
}

func SetCookie(h CookieHandler, w http.ResponseWriter, ticket *auth.Ticket, cookie *Cookie) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,

		Domain:  string(h.Domain()),
		Path:    "/",
		Expires: ticket.Expires(),

		Secure:   true,
		HttpOnly: true,
	})
}

func SetAuthCookie(h CredentialHandler, w http.ResponseWriter, ticket *auth.Ticket) error {
	err := SetCredentialCookie(h, w, ticket)
	if err != nil {
		return err
	}

	err = SetCloudFrontCookie(h, w, ticket)
	if err != nil {
		return err
	}

	return nil
}
