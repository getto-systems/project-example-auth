package authenticate

import (
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrPasswordMatchFailed = errors.New("password match failed")
	ErrTicketIssueFailed   = errors.New("ticket issue failed")
)

type PasswordAuthenticator struct {
	passwordRepository PasswordRepository
	issuerRepository   IssuerRepository

	user    user.User
	request data.Request
}

func (authenticator PasswordAuthenticator) MatchPassword(password data.RawPassword) (data.Token, error) {
	authenticator.user.PasswordMatching(authenticator.request)

	passwordMatcher := authenticator.passwordRepository.Find(authenticator.user)

	err := passwordMatcher.Match(password)
	if err != nil {
		authenticator.user.PasswordMatchFailed(authenticator.request, err)
		return nil, ErrPasswordMatchFailed
	}

	issuer := authenticator.issuerRepository.Find(authenticator.user)
	ticket, token, err := issuer.Authenticated(authenticator.request.RequestedAt)
	if err != nil {
		authenticator.user.TicketIssueFailed(authenticator.request, err)
		return nil, ErrTicketIssueFailed
	}

	user := authenticator.user.Authenticated(ticket)
	user.Authenticated(authenticator.request)

	return token, nil
}

type PasswordAuthenticatorFactory struct {
	passwordRepository PasswordRepository
	issuerRepository   IssuerRepository
	userFactory        user.UserFactory
}

func NewPasswordAuthenticatorFactory(passwordRepository PasswordRepository, issuerRepository IssuerRepository, userFactory user.UserFactory) PasswordAuthenticatorFactory {
	return PasswordAuthenticatorFactory{
		passwordRepository: passwordRepository,
		issuerRepository:   issuerRepository,
		userFactory:        userFactory,
	}
}

func (f PasswordAuthenticatorFactory) New(userID data.UserID, request data.Request) PasswordAuthenticator {
	return PasswordAuthenticator{
		passwordRepository: f.passwordRepository,
		issuerRepository:   f.issuerRepository,

		user:    f.userFactory.New(userID),
		request: request,
	}
}
