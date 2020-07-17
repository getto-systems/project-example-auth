package password

import (
	"errors"

	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

var (
	ErrValidateFailed = errors.New("password-validate-failed")
	ErrRegisterFailed = errors.New("password-register-failed")
)

type EventPublisher interface {
	registerEventPublisher
	validateEventPublisher
}

type EventHandler interface {
	EventPublisher
}

type DB interface {
	registerDB
	validateDB
}

type PasswordValidater struct {
	validater   Validater
	issuer      ticket.TicketIssuer
	tokenIssuer ticket.TokenIssuer
}

func NewPasswordValidater(
	validater Validater,
	issuer ticket.TicketIssuer,
	api ticket.ApiTokenIssuer,
	content ticket.ContentTokenIssuer,
) PasswordValidater {
	return PasswordValidater{
		validater:   validater,
		issuer:      issuer,
		tokenIssuer: ticket.NewTokenIssuer(api, content),
	}
}

func (usecase PasswordValidater) Validate(request data.Request, user data.User, password data.RawPassword) (ticket.Ticket, ticket.Nonce, ticket.ApiToken, ticket.ContentToken, data.Expires, error) {
	err := usecase.validater.validate(request, user, password)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, ErrValidateFailed
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
	ticketValidater   ticket.TicketValidater
	passwordValidater Validater
	register          Register
}

func NewPasswordRegister(
	ticketValidater ticket.TicketValidater,
	passwordValidater Validater,
	register Register,
) PasswordRegister {
	return PasswordRegister{
		ticketValidater:   ticketValidater,
		passwordValidater: passwordValidater,
		register:          register,
	}
}

type PasswordRegisterParam struct {
	OldPassword data.RawPassword
	NewPassword data.RawPassword
}

func (usecase PasswordRegister) Register(request data.Request, ticket ticket.Ticket, nonce ticket.Nonce, password PasswordRegisterParam) error {
	user, err := usecase.ticketValidater.Validate(request, ticket, nonce)
	if err != nil {
		return ErrRegisterFailed
	}

	err = usecase.passwordValidater.validate(request, user, password.OldPassword)
	if err != nil {
		return ErrRegisterFailed
	}

	err = usecase.register.register(request, user, password.NewPassword)
	if err != nil {
		return ErrRegisterFailed
	}

	return nil
}
