package renew

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/handler"
	"github.com/getto-systems/project-example-id/signature"
)

type Handler struct {
	domain     handler.Domain
	tokener    auth.Tokener
	cloudfront signature.CloudFrontSigner

	userRepository auth.UserRepository
}

func (h Handler) Domain() handler.Domain {
	return h.domain
}

func (h Handler) Tokener() auth.Tokener {
	return h.tokener
}

func (h Handler) CloudFrontSigner() signature.CloudFrontSigner {
	return h.cloudfront
}

func NewHandler(domain handler.Domain, tokener auth.Tokener, cloudfront signature.CloudFrontSigner, userRepository auth.UserRepository) Handler {
	return Handler{
		domain,
		tokener,
		cloudfront,
		userRepository,
	}
}

type Input struct {
	Path string `json:"path"`
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := auth.Path(input.Path)

	ticket, err := handler.ParseCookieToken(h, r, path)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now().UTC()

	if ticket.IsRenewRequired(now) {
		user, err := auth.NewUserFromRepository(h.userRepository, ticket.User().UserID(), path)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ticket = user.NewTicket(now)

		err = handler.SetAuthCookie(h, w, ticket)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	response, err := h.tokener.FullToken(ticket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", response)
}
