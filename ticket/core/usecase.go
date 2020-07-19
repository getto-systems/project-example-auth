package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type usecase struct {
	issuer    issuer
	extender  extender
	validator validator
	shrinker  shrinker

	apiTokenIssuer     apiTokenIssuer
	contentTokenIssuer contentTokenIssuer
}

func NewUsecase(
	pub ticket.EventPublisher,
	db ticket.DB,
	exp ticket.Expiration,
	signer ticket.Signer,
	apiTokenSigner ticket.ApiTokenSigner,
	contentTokenSigner ticket.ContentTokenSigner,
	gen ticket.NonceGenerator,
) ticket.Usecase {
	return usecase{
		issuer:    newIssuer(pub, db, signer, exp, gen),
		extender:  newExtender(pub, db, signer, exp),
		validator: newValidator(pub, signer),
		shrinker:  newShrinker(pub, db),

		apiTokenIssuer:     newApiTokenIssuer(pub, db, apiTokenSigner),
		contentTokenIssuer: newContentTokenIssuer(pub, contentTokenSigner),
	}
}

func (usecase usecase) Issue(request data.Request, user data.User) (ticket.Ticket, ticket.Nonce, data.Expires, error) {
	newTicket, nonce, expires, err := usecase.issuer.issue(request, user)
	if err != nil {
		return nil, "", data.Expires{}, ticket.ErrIssueFailed
	}

	return newTicket, nonce, expires, nil
}

func (usecase usecase) Extend(request data.Request, originalTicket ticket.Ticket, nonce ticket.Nonce) (
	ticket.Ticket,
	ticket.ApiToken,
	ticket.ContentToken,
	data.Expires,
	error,
) {
	user, err := usecase.validator.validate(request, originalTicket, nonce)
	if err != nil {
		return nil, nil, nil, data.Expires{}, ticket.ErrValidateFailed
	}

	extendedTicket, expires, err := usecase.extender.extend(request, nonce, user)
	if err != nil {
		return nil, nil, nil, data.Expires{}, ticket.ErrExtendFailed
	}

	apiToken, err := usecase.apiTokenIssuer.issue(request, user, expires)
	if err != nil {
		return nil, nil, nil, data.Expires{}, err
	}

	contentToken, err := usecase.contentTokenIssuer.issue(request, user, expires)
	if err != nil {
		return nil, nil, nil, data.Expires{}, err
	}

	return extendedTicket, apiToken, contentToken, expires, nil
}

func (usecase usecase) Validate(request data.Request, originalTicket ticket.Ticket, nonce ticket.Nonce) (data.User, error) {
	user, err := usecase.validator.validate(request, originalTicket, nonce)
	if err != nil {
		return data.User{}, ticket.ErrValidateFailed
	}
	return user, nil
}

func (usecase usecase) Shrink(request data.Request, originalTicket ticket.Ticket, nonce ticket.Nonce) error {
	user, err := usecase.validator.validate(request, originalTicket, nonce)
	if err != nil {
		return ticket.ErrValidateFailed
	}

	err = usecase.shrinker.shrink(request, nonce, user)
	if err != nil {
		return ticket.ErrShrinkFailed
	}

	return nil
}

func (usecase usecase) IssueToken(request data.Request, user data.User, expires data.Expires) (
	ticket.ApiToken,
	ticket.ContentToken,
	error,
) {
	apiToken, err := usecase.apiTokenIssuer.issue(request, user, expires)
	if err != nil {
		return nil, nil, ticket.ErrApiTokenIssueFailed
	}

	contentToken, err := usecase.contentTokenIssuer.issue(request, user, expires)
	if err != nil {
		return nil, nil, ticket.ErrContentTokenIssueFailed
	}

	return apiToken, contentToken, nil
}
