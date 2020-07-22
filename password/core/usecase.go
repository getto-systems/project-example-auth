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
	pub password.EventPublisher,
	db password.DB,

	matcher password.Matcher,
	gen password.Generator,
	resetGen password.ResetGenerator,

	ticket ticket.Usecase,
) password.Usecase {
	return usecase{
		validator:  newValidator(pub, db, matcher),
		registerer: newRegisterer(pub, db, gen),
		resetter:   newResetter(pub, db, resetGen),

		ticket: ticket,
	}
}

func (usecase usecase) Validate(request data.Request, login password.Login, rawPassword password.RawPassword) (
	ticket.Ticket,
	ticket.Nonce,
	ticket.ApiToken,
	ticket.ContentToken,
	data.Expires,
	error,
) {
	user, err := usecase.validator.validate(request, login, rawPassword)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	newTicket, nonce, expires, err := usecase.ticket.Issue(request, user)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	apiToken, contentToken, err := usecase.ticket.IssueToken(request, user, expires)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	return newTicket, nonce, apiToken, contentToken, expires, nil
}

func (usecase usecase) GetLogin(request data.Request, ticket ticket.Ticket, nonce ticket.Nonce) (password.Login, error) {
	user, err := usecase.ticket.Validate(request, ticket, nonce)
	if err != nil {
		return password.Login{}, err
	}

	login, err := usecase.registerer.getLogin(request, user)
	if err != nil {
		return password.Login{}, password.ErrLoginNotFound
	}

	return login, nil
}

func (usecase usecase) Register(
	request data.Request,
	ticket ticket.Ticket,
	nonce ticket.Nonce,
	login password.Login,
	param password.RegisterParam,
) error {
	user, err := usecase.ticket.Validate(request, ticket, nonce)
	if err != nil {
		return err
	}

	user, err = usecase.validator.validate(request, login, param.OldPassword)
	if err != nil {
		return err
	}

	err = usecase.registerer.register(request, user, param.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

func (usecase usecase) IssueReset(request data.Request, login password.Login) (password.Reset, error) {
	reset, err := usecase.resetter.issueReset(request, login)
	if err != nil {
		return password.Reset{}, err
	}
	return reset, nil
}

func (usecase usecase) GetResetStatus(request data.Request, reset password.Reset) (password.ResetStatus, error) {
	status, err := usecase.resetter.getResetStatus(request, reset)
	if err != nil {
		return password.ResetStatus{}, err
	}
	return status, nil
}

func (usecase usecase) Reset(
	request data.Request,
	login password.Login,
	token password.ResetToken,
	newPassword password.RawPassword,
) (
	ticket.Ticket,
	ticket.Nonce,
	ticket.ApiToken,
	ticket.ContentToken,
	data.Expires,
	error,
) {
	user, err := usecase.resetter.validate(request, login, token)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	err = usecase.registerer.register(request, user, newPassword)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	newTicket, nonce, expires, err := usecase.ticket.Issue(request, user)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	apiToken, contentToken, err := usecase.ticket.IssueToken(request, user, expires)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, err
	}

	return newTicket, nonce, apiToken, contentToken, expires, nil
}
