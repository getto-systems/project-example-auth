package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/authenticate"

	"github.com/getto-systems/project-example-id/data"

	"fmt"
)

const COOKIE_SIGNED_TICKET = "Getto-Example-Auth-Signed-Ticket"

type AuthResponse struct {
	cookieDomain CookieDomain
	issuer       Issuer
}

type CookieDomain string

type Issuer struct {
	awsCloudFront AwsCloudFrontIssuer
	app           AppIssuer
}

func NewAuthResponse(cookieDomain CookieDomain, awsCloudFront AwsCloudFrontIssuer, app AppIssuer) AuthResponse {
	return AuthResponse{
		cookieDomain: cookieDomain,
		issuer: Issuer{
			awsCloudFront: awsCloudFront,
			app:           app,
		},
	}
}

func SignedTicket(r *http.Request) (data.SignedTicket, error) {
	cookie, err := r.Cookie(COOKIE_SIGNED_TICKET)
	if err != nil {
		return nil, err
	}

	return data.SignedTicket(cookie.Value), nil
}

func (response AuthResponse) write(w http.ResponseWriter, ticket data.Ticket, signedTicket data.SignedTicket, logger http_handler.Logger) {
	awsCloudFrontTicket, err := response.issuer.awsCloudFront.Authorized(ticket)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	appToken, err := response.issuer.app.Authorized(ticket)
	if err != nil {
		errorResponse(w, err, logger)
		return
	}

	response.setSignedTicketCookie(w, signedTicket, ticket.Expires)
	response.setAwsCloudFrontCookie(w, awsCloudFrontTicket, ticket.Expires)

	jsonResponse(w, appToken)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", jsonData)
}

func errorResponse(w http.ResponseWriter, err error, logger http_handler.Logger) {
	w.Header().Set("Content-Type", "application/json")

	switch err {

	case
		http_handler.ErrEmptyBody,
		http_handler.ErrBodyParseFailed:

		w.WriteHeader(http.StatusBadRequest)

	case
		ErrSignedTicketCookieNotFound,
		authenticate.ErrTicketAuthFailed,
		authenticate.ErrPasswordAuthFailed:

		w.WriteHeader(http.StatusUnauthorized)

	default:
		logger.DebugError("internal server error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (response AuthResponse) setSignedTicketCookie(w http.ResponseWriter, signedTicket data.SignedTicket, expires data.Expires) {
	http.SetCookie(w, &http.Cookie{
		Name:  COOKIE_SIGNED_TICKET,
		Value: string(signedTicket),

		Domain:  string(response.cookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (response AuthResponse) setAwsCloudFrontCookie(w http.ResponseWriter, ticket AwsCloudFrontTicket, expires data.Expires) {
	http.SetCookie(w, &http.Cookie{
		Name:  "CloudFront-Key-Pair-Id",
		Value: string(ticket.KeyPairID),

		Domain:  string(response.cookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "CloudFront-Policy",
		Value: string(ticket.Token.Policy),

		Domain:  string(response.cookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "CloudFront-Signature",
		Value: string(ticket.Token.Signature),

		Domain:  string(response.cookieDomain),
		Path:    "/",
		Expires: expires.Time(),

		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
