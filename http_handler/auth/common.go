package auth

import (
	"encoding/json"
	"net/http"

	"github.com/getto-systems/project-example-id/auth"

	"errors"
	"fmt"
)

var ErrBodyNotSent = errors.New("body not sent")
var ErrBodyParseFailed = errors.New("body parse failed")
var ErrTicketCookieNotSent = errors.New("ticket cookie not sent")
var ErrResponseEncodeFailed = errors.New("response encode failed")

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
