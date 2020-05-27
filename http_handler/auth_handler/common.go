package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"errors"
	"fmt"
)

const COOKIE_AUTH_TOKEN = "Getto-Example-Auth-Token"

type CookieDomain string

type Cookie struct {
	Name  string
	Value string
}

type CookieSetter struct {
	ResponseWriter http.ResponseWriter
	CookieDomain   CookieDomain
	Ticket         user.Ticket
}

var (
	ErrBodyNotSent          = errors.New("body not sent")
	ErrBodyParseFailed      = errors.New("body parse failed")
	ErrTicketCookieNotSent  = errors.New("ticket cookie not sent")
	ErrResponseEncodeFailed = errors.New("response encode failed")
)

func jsonResponse(w http.ResponseWriter, response interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}

func httpStatusCode(err error) int {
	switch err {

	case ErrBodyNotSent, ErrBodyParseFailed:
		return http.StatusBadRequest

	case ErrTicketCookieNotSent, auth.ErrRenewTokenParseFailed:
		return http.StatusUnauthorized

	case auth.ErrUserAccessDenied:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}

func setAuthTokenCookie(w http.ResponseWriter, cookieDomain CookieDomain, ticket user.Ticket, token auth.Token) {
	setter := CookieSetter{
		ResponseWriter: w,
		CookieDomain:   cookieDomain,
		Ticket:         ticket,
	}
	setter.setTicketCookie(token.RenewToken)
	setter.setAwsCloudFrontCookie(token.AwsCloudFrontToken)
}

func getRenewToken(r *http.Request) (token.RenewToken, error) {
	cookie, err := r.Cookie(COOKIE_AUTH_TOKEN)
	if err != nil {
		return nil, err
	}

	return token.RenewToken(cookie.Value), nil
}

func (setter CookieSetter) setTicketCookie(renewToken token.RenewToken) {
	setter.setCookie(&Cookie{
		Name:  COOKIE_AUTH_TOKEN,
		Value: string(renewToken),
	})
}

func (setter CookieSetter) setAwsCloudFrontCookie(token token.AwsCloudFrontToken) {
	setter.setCookie(&Cookie{
		Name:  "CloudFront-Key-Pair-Id",
		Value: string(token.KeyPairID),
	})

	setter.setCookie(&Cookie{
		Name:  "CloudFront-Policy",
		Value: string(token.Policy),
	})

	setter.setCookie(&Cookie{
		Name:  "CloudFront-Signature",
		Value: string(token.Signature),
	})
}

func (setter CookieSetter) setCookie(cookie *Cookie) {
	http.SetCookie(setter.ResponseWriter, &http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,

		Domain:  string(setter.CookieDomain),
		Path:    "/",
		Expires: setter.Ticket.Expires(),

		Secure:   true,
		HttpOnly: true,
	})
}
