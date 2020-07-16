package password

import (
	"errors"

	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

var (
	ErrVerifyFailed   = errors.New("password-verify-failed")
	ErrRegisterFailed = errors.New("password-register-failed")
)

type EventPublisher interface {
	ValidatePassword(data.Request, data.User)
	ValidatePasswordFailed(data.Request, data.User, error)
	PasswordRegistered(data.Request, data.User)

	VerifyPassword(data.Request, data.User)
	VerifyPasswordFailed(data.Request, data.User, error)
	AuthenticatedByPassword(data.Request, data.User)
}

type EventHandler interface {
	EventPublisher
}

type DB interface {
	RegisterUserPassword(data.User, data.HashedPassword) error

	FindUserPassword(data.User) (data.HashedPassword, error)
}

type PasswordVerifier struct {
	verifier    Verifier
	issuer      ticket.TicketIssuer
	tokenIssuer ticket.TokenIssuer
}

func NewPasswordVerifier(
	verifier Verifier,
	issuer ticket.TicketIssuer,
	api ticket.ApiTokenIssuer,
	content ticket.ContentTokenIssuer,
) PasswordVerifier {
	return PasswordVerifier{
		verifier:    verifier,
		issuer:      issuer,
		tokenIssuer: ticket.NewTokenIssuer(api, content),
	}
}

func (usecase PasswordVerifier) Verify(request data.Request, user data.User, password data.RawPassword) (ticket.Ticket, ticket.Nonce, ticket.ApiToken, ticket.ContentToken, data.Expires, error) {
	err := usecase.verifier.verify(request, user, password)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, ErrVerifyFailed
	}

	newTicket, nonce, expires, err := usecase.issuer.Issue(request, user)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	apiToken, contentToken, err := usecase.tokenIssuer.Issue(request, user, expires)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	return newTicket, nonce, apiToken, contentToken, expires, nil
}

type PasswordRegister struct {
	ticketVerifier   ticket.TicketVerifier
	passwordVerifier Verifier
	register         Register
}

func NewPasswordRegister(
	ticketVerifier ticket.TicketVerifier,
	passwordVerifier Verifier,
	register Register,
) PasswordRegister {
	return PasswordRegister{
		ticketVerifier:   ticketVerifier,
		passwordVerifier: passwordVerifier,
		register:         register,
	}
}

type PasswordRegisterParam struct {
	OldPassword data.RawPassword
	NewPassword data.RawPassword
}

func (usecase PasswordRegister) Register(request data.Request, ticket ticket.Ticket, nonce ticket.Nonce, password PasswordRegisterParam) error {
	user, err := usecase.ticketVerifier.Verify(request, ticket, nonce)
	if err != nil {
		return ErrRegisterFailed
	}

	err = usecase.passwordVerifier.verify(request, user, password.OldPassword)
	if err != nil {
		return ErrRegisterFailed
	}

	err = usecase.register.register(request, user, password.NewPassword)
	if err != nil {
		return ErrRegisterFailed
	}

	return nil
}
