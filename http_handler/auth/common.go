package auth

import (
	"errors"
	"net/http"

	"github.com/getto-systems/project-example-id/auth"
)

var ErrBodyNotSent = errors.New("body not sent")
var ErrBodyParseFailed = errors.New("body parse failed")
var ErrTicketCookieNotSent = errors.New("ticket cookie not sent")

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
