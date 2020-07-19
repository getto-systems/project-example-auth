package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"
)

type usecase struct {
	validator  validator
	registerer registerer

	ticket ticket.Usecase
}

func NewUsecase(
	pub password.EventPublisher,
	db password.DB,

	matcher password.Matcher,
	gen password.Generator,

	ticket ticket.Usecase,
) password.Usecase {
	return usecase{
		validator:  newValidator(pub, db, matcher),
		registerer: newRegisterer(pub, db, gen),

		ticket: ticket,
	}
}

func (usecase usecase) Validate(request data.Request, user data.User, rawPassword password.RawPassword) (
	ticket.Ticket,
	ticket.Nonce,
	ticket.ApiToken,
	ticket.ContentToken,
	data.Expires,
	error,
) {
	err := usecase.validator.validate(request, user, rawPassword)
	if err != nil {
		return nil, "", nil, nil, data.Expires{}, password.ErrValidateFailed
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

func (usecase usecase) Register(
	request data.Request,
	ticket ticket.Ticket,
	nonce ticket.Nonce,
	param password.RegisterParam,
) error {
	user, err := usecase.ticket.Validate(request, ticket, nonce)
	if err != nil {
		return password.ErrRegisterFailed
	}

	err = usecase.validator.validate(request, user, param.OldPassword)
	if err != nil {
		return password.ErrRegisterFailed
	}

	err = usecase.registerer.register(request, user, param.NewPassword)
	if err != nil {
		return password.ErrRegisterFailed
	}

	return nil
}
