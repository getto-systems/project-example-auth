package auth

import (
	"errors"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/ticket"
	"github.com/getto-systems/project-example-auth/user"
)

var (
	ErrBadRequest    = errors.New("bad-request")
	ErrInvalidTicket = errors.New("invalid-ticket")
	ErrServerError   = errors.New("server-error")
)

type (
	Backend struct {
		ticket        ticket.Action
		credential    credential.Action
		user          user.Action
		password      password.Action
		passwordReset password_reset.Action
	}

	Usecase struct {
		Backend
		handler CredentialHandler
	}

	CredentialHandler interface {
		GetTicketNonceAndSignature() (credential.TicketNonce, credential.TicketSignature, error)
		SetCredential(credential.Credential)
		ClearCredential()
	}
)

func NewBackend(
	ticket ticket.Action,
	credential credential.Action,
	user user.Action,
	password password.Action,
	passwordReset password_reset.Action,
) Backend {
	return Backend{
		ticket:        ticket,
		credential:    credential,
		user:          user,
		password:      password,
		passwordReset: passwordReset,
	}
}

func NewUsecase(backend Backend, handler CredentialHandler) Usecase {
	return Usecase{
		Backend: backend,
		handler: handler,
	}
}

func (u Usecase) validateTicket(request request.Request) (_ user.User, _ credential.TicketNonce, err error) {
	nonce, signature, err := u.handler.GetTicketNonceAndSignature()
	if err != nil {
		switch err {
		default:
			err = ErrInvalidTicket
		}
		return
	}

	user, err := u.credential.ParseTicketSignature(request, nonce, signature)
	if err != nil {
		switch err {
		case credential.ErrParseTicketSignatureParseFailed,
			credential.ErrParseTicketSignatureMatchFailedNonce:

			err = ErrInvalidTicket
			u.clearCredential()

		default:
			err = ErrServerError
		}
		return
	}

	err = u.ticket.Validate(request, user, nonce)
	if err != nil {
		switch err {
		case ticket.ErrValidateAlreadyExpired,
			ticket.ErrValidateMatchFailedUser,
			ticket.ErrValidateNotFoundTicket:

			err = ErrInvalidTicket
			u.clearCredential()

		default:
			err = ErrServerError
		}
		return
	}

	return user, nonce, nil
}

func (u Usecase) issueCredential(request request.Request, ticket credential.Ticket) (err error) {
	ticketToken, err := u.credential.IssueTicketToken(request, ticket)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	apiToken, err := u.credential.IssueApiToken(request, ticket)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	contentToken, err := u.credential.IssueContentToken(request, ticket)
	if err != nil {
		switch err {
		default:
			err = ErrServerError
		}
		return
	}

	u.handler.SetCredential(credential.NewCredential(ticketToken, apiToken, contentToken))

	return nil
}
func (u Usecase) clearCredential() {
	u.handler.ClearCredential()
}
