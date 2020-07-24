package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"
)

type usecase struct {
	validator  validator
	registerer registerer
	resetter   resetter

	ticket ticket.Usecase
}

func NewUsecase(
	logger password.Logger,

	passwords password.PasswordRepository,
	sessions password.ResetSessionRepository,

	exp password.ResetSessionExpiration,
	encrypter password.PasswordEncrypter,
	gen password.ResetSessionGenerator,

	ticket ticket.Usecase,
) password.Usecase {
	return usecase{
		validator:  newValidator(logger, passwords, encrypter),
		registerer: newRegisterer(logger, passwords, encrypter),
		resetter:   newResetter(logger, passwords, sessions, exp, gen),

		ticket: ticket,
	}
}

func (usecase usecase) Validate(request data.Request, login password.Login, rawPassword password.RawPassword) (_ ticket.Ticket, _ ticket.Nonce, _ ticket.ApiToken, _ ticket.ContentToken, _ data.Expires, err error) {
	user, err := usecase.validator.validate(request, login, rawPassword)
	if err != nil {
		return
	}

	return usecase.issueTicket(request, user)
}

func (usecase usecase) GetLogin(request data.Request, ticket ticket.Ticket, nonce ticket.Nonce) (_ password.Login, err error) {
	user, err := usecase.ticket.Validate(request, ticket, nonce)
	if err != nil {
		return
	}

	login, err := usecase.registerer.getLogin(request, user)
	if err != nil {
		return
	}

	return login, nil
}

func (usecase usecase) Register(request data.Request, ticket ticket.Ticket, nonce ticket.Nonce, login password.Login, param password.RegisterParam) (err error) {
	user, err := usecase.ticket.Validate(request, ticket, nonce)
	if err != nil {
		return
	}

	user, err = usecase.validator.validate(request, login, param.OldPassword)
	if err != nil {
		return
	}

	err = usecase.registerer.register(request, user, param.NewPassword)
	if err != nil {
		return
	}

	return nil
}

func (usecase usecase) CreateResetSession(request data.Request, login password.Login) (password.ResetSession, error) {
	return usecase.resetter.createResetSession(request, login)
}

func (usecase usecase) GetResetStatus(request data.Request, session password.ResetSession) (password.ResetStatus, error) {
	return usecase.resetter.getResetStatus(request, session)
}

func (usecase usecase) Reset(request data.Request, login password.Login, token password.ResetToken, newPassword password.RawPassword) (_ ticket.Ticket, _ ticket.Nonce, _ ticket.ApiToken, _ ticket.ContentToken, _ data.Expires, err error) {
	user, err := usecase.resetter.validate(request, login, token)
	if err != nil {
		return
	}

	err = usecase.registerer.register(request, user, newPassword)
	if err != nil {
		return
	}

	return usecase.issueTicket(request, user)
}

func (usecase usecase) issueTicket(request data.Request, user data.User) (_ ticket.Ticket, _ ticket.Nonce, _ ticket.ApiToken, _ ticket.ContentToken, _ data.Expires, err error) {
	ticket, nonce, expires, err := usecase.ticket.Issue(request, user)
	if err != nil {
		return
	}

	api, content, err := usecase.ticket.IssueToken(request, user, expires)
	if err != nil {
		return
	}

	return ticket, nonce, api, content, expires, nil
}
