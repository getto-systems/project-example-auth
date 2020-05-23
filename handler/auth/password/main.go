package password

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/auth/password"
	"github.com/getto-systems/project-example-id/handler"
	"github.com/getto-systems/project-example-id/signature"
)

type Handler struct {
	domain     handler.Domain
	tokener    auth.Tokener
	cloudfront signature.CloudFrontSigner

	userRepository         auth.UserRepository
	userPasswordRepository password.UserPasswordRepository
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

func NewHandler(domain handler.Domain, tokener auth.Tokener, cloudfront signature.CloudFrontSigner, userRepository auth.UserRepository, userPasswordRepository password.UserPasswordRepository) Handler {
	return Handler{
		domain,
		tokener,
		cloudfront,
		userRepository,
		userPasswordRepository,
	}
}

type Input struct {
	Path     string `json:"path"`
	UserID   string `json:"user_id"`
	Password string `json:"password"`
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

	userID := auth.UserID(input.UserID)
	userPassword := password.NewUserPassword(h.userPasswordRepository, userID)
	path := auth.Path(input.Path)

	if !userPassword.Match(password.Password(input.Password)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := auth.NewUserFromRepository(h.userRepository, userID, path)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ticket := user.NewTicket(time.Now().UTC())

	err = handler.SetAuthCookie(h, w, ticket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := h.tokener.FullToken(ticket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", response)
}
