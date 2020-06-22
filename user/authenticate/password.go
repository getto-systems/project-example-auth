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
	authenticator.passwordMatching()

	passwordMatcher := authenticator.passwordMatcher()

	err := passwordMatcher.Match(password)
	if err != nil {
		authenticator.passwordMatchFailed(err)
		return nil, ErrPasswordMatchFailed
	}

	issuer := authenticator.issuer()
	token, err := issuer.Authenticated(authenticator.request.RequestedAt)
	if err != nil {
		authenticator.ticketIssueFailed(err)
		return nil, ErrTicketIssueFailed
	}

	authenticator.authenticated(issuer.Profile())

	return token, nil
}

func (authenticator PasswordAuthenticator) passwordMatcher() user.PasswordMatcher {
	return authenticator.passwordRepository.Find(authenticator.user.UserID())
}

func (authenticator PasswordAuthenticator) issuer() user.Issuer {
	return authenticator.issuerRepository.Find(authenticator.user.UserID())
}

func (authenticator PasswordAuthenticator) passwordMatching() {
	authenticator.user.PasswordMatching(authenticator.request)
}

func (authenticator PasswordAuthenticator) passwordMatchFailed(err error) {
	authenticator.user.PasswordMatchFailed(authenticator.request, err)
}

func (authenticator PasswordAuthenticator) ticketIssueFailed(err error) {
	authenticator.user.TicketIssueFailed(authenticator.request, err)
}

func (authenticator PasswordAuthenticator) authenticated(profile data.Profile) {
	authenticator.user.Authenticated(authenticator.request, profile)
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
