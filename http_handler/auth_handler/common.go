package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/applog"

	"github.com/getto-systems/project-example-id/basic"
	"github.com/getto-systems/project-example-id/http_handler/auth_handler/token"

	"errors"
	"fmt"
	"time"
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
	Expires        basic.Expires
}

var (
	ErrBodyNotSent                       = errors.New("body not sent")
	ErrBodyParseFailed                   = errors.New("body parse failed")
	ErrTicketCookieNotSent               = errors.New("ticket cookie not sent")
	ErrResponseEncodeFailed              = errors.New("response encode failed")
	ErrRenewTokenParseFailed             = errors.New("ticket token parse failed")
	ErrRenewTokenSerializeFailed         = errors.New("renew token serialize failed")
	ErrAwsCloudFrontTokenSerializeFailed = errors.New("aws cloudfront token serialize failed")
	ErrAppTokenSerializeFailed           = errors.New("app token serialize failed")
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

	case ErrTicketCookieNotSent, ErrRenewTokenParseFailed:
		return http.StatusUnauthorized

	case auth.ErrUserAccessDenied:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}

func setAuthTokenCookie(w http.ResponseWriter, cookieDomain CookieDomain, token Token) {
	setter := CookieSetter{
		ResponseWriter: w,
		CookieDomain:   cookieDomain,
		Expires:        token.Expires,
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
		Expires: time.Time(setter.Expires),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

type Token struct {
	Expires basic.Expires

	RenewToken         token.RenewToken
	AwsCloudFrontToken token.AwsCloudFrontToken
}

func (token Token) String() string {
	return fmt.Sprintf(
		"Token{Expires: %s, RenewToken:%s, AwsCloudFrontToken:%s}",
		token.Expires,
		token.RenewToken,
		token.AwsCloudFrontToken,
	)
}

type AuthHandler interface {
	Logger() applog.Logger
	TicketSerializer() token.TicketSerializer
	AwsCloudFrontSerializer() token.AwsCloudFrontSerializer
}

func SerializeAuthToken(handler AuthHandler, ticket basic.Ticket) (Token, error) {
	logger := handler.Logger()

	logger.Debugf("serialize ticket: %v", ticket)

	logger.Debug("serialize renew token...")

	renewToken, err := handler.TicketSerializer().RenewToken(ticket)
	if err != nil {
		logger.Errorf("ticket serialize error: %s; %v", err, ticket)
		return Token{}, ErrRenewTokenSerializeFailed
	}

	logger.Debug("serialize aws cloudfront token...")

	awsCloudFrontToken, err := handler.AwsCloudFrontSerializer().Token(ticket)
	if err != nil {
		logger.Errorf("aws cloudfront serialize error: %s; %v", err, ticket)
		return Token{}, ErrAwsCloudFrontTokenSerializeFailed
	}

	logger.Debug("handling ticket token...")

	return Token{
		Expires: ticket.Expires,

		RenewToken:         renewToken,
		AwsCloudFrontToken: awsCloudFrontToken,
	}, nil
}

func SerializeAppToken(handler AuthHandler, ticket basic.Ticket) (token.AppToken, error) {
	logger := handler.Logger()

	logger.Debugf("serialize ticket for app: %v", ticket)

	appToken, err := handler.TicketSerializer().AppToken(ticket)
	if err != nil {
		logger.Errorf("ticket serialize error: %s; %v", err, ticket)
		return token.AppToken{}, ErrAppTokenSerializeFailed
	}

	return appToken, nil
}
