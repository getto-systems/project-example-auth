package authenticate

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"
)

type AuthByPassword struct {
	authenticated user.UserAuthenticatedFactory
	passwordAuth  user.UserPasswordAuthFactory
}

func (auth AuthByPassword) Authenticate(request data.Request, userID data.UserID, password data.RawPassword) (data.Ticket, data.SignedTicket, error) {
	passwordAuthUser := auth.passwordAuth.New(request, userID)
	err := passwordAuthUser.Authenticate(password)
	if err != nil {
		return data.Ticket{}, nil, ErrPasswordAuthFailed
	}

	authenticatedUser := auth.authenticated.New(request, userID)
	return authenticatedUser.IssueTicket()
}

func NewAuthByPassword(authenticated user.UserAuthenticatedFactory, passwordAuth user.UserPasswordAuthFactory) AuthByPassword {
	return AuthByPassword{
		authenticated: authenticated,
		passwordAuth:  passwordAuth,
	}
}
