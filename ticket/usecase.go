package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

var (
	ErrValidateFailed          = errors.New("ticket-validate-failed")
	ErrExtendFailed            = errors.New("ticket-extend-failed")
	ErrShrinkFailed            = errors.New("ticket-shrink-failed")
	ErrIssueFailed             = errors.New("ticket-issue-failed")
	ErrApiTokenIssueFailed     = errors.New("api-token-issue-failed")
	ErrContentTokenIssueFailed = errors.New("content-token-issue-failed")
)

type (
	Ticket       []byte
	ApiToken     []byte
	ContentToken interface {
		Policy() string
		Signature() string
	}

	Nonce string

	Expiration struct {
		expires     data.Second
		extendLimit data.Second
	}
)

type EventPublisher interface {
	apiTokenIssueEventPublisher
	contentTokenIssueEventPublisher
	extendEventPublisher
	issueEventPublisher
	shrinkEventPublisher
	validateEventPublisher
}

type EventHandler interface {
	EventPublisher
}

type DB interface {
	apiTokenDB
	extendDB
	issueDB
	shrinkDB
}

type TicketExtender struct {
	validater   Validater
	extender    Extender
	tokenIssuer TokenIssuer
}

func NewTicketExtender(
	validater Validater,
	extender Extender,
	api ApiTokenIssuer,
	content ContentTokenIssuer,
) TicketExtender {
	return TicketExtender{
		validater:   validater,
		extender:    extender,
		tokenIssuer: NewTokenIssuer(api, content),
	}
}

func (usecase TicketExtender) Extend(request data.Request, ticket Ticket, nonce Nonce) (Ticket, ApiToken, ContentToken, data.Expires, error) {
	user, err := usecase.validater.validate(request, ticket, nonce)
	if err != nil {
		return nil, nil, nil, data.Expires{}, ErrValidateFailed
	}

	ticket, expires, err := usecase.extender.extend(request, nonce, user)
	if err != nil {
		return nil, nil, nil, data.Expires{}, ErrExtendFailed
	}

	apiToken, contentToken, err := usecase.tokenIssuer.Issue(request, user, expires)
	if err != nil {
		return nil, nil, nil, data.Expires{}, err
	}

	return ticket, apiToken, contentToken, expires, nil
}

type TicketValidater struct {
	validater Validater
}

func NewTicketValidater(validater Validater) TicketValidater {
	return TicketValidater{
		validater: validater,
	}
}

func (usecase TicketValidater) Validate(request data.Request, ticket Ticket, nonce Nonce) (data.User, error) {
	user, err := usecase.validater.validate(request, ticket, nonce)
	if err != nil {
		return data.User{}, ErrValidateFailed
	}
	return user, nil
}

type TicketShrinker struct {
	validater Validater
	shrinker  Shrinker
}

func NewTicketShrinker(
	validater Validater,
	shrinker Shrinker,
) TicketShrinker {
	return TicketShrinker{
		validater: validater,
		shrinker:  shrinker,
	}
}

func (usecase TicketShrinker) Shrink(request data.Request, ticket Ticket, nonce Nonce) error {
	user, err := usecase.validater.validate(request, ticket, nonce)
	if err != nil {
		return ErrValidateFailed
	}

	err = usecase.shrinker.shrink(request, nonce, user)
	if err != nil {
		return ErrShrinkFailed
	}

	return nil
}

type TicketIssuer struct {
	issuer Issuer
}

func NewTicketIssuer(issuer Issuer) TicketIssuer {
	return TicketIssuer{
		issuer: issuer,
	}
}

func (usecase TicketIssuer) Issue(request data.Request, user data.User) (Ticket, Nonce, data.Expires, error) {
	ticket, nonce, expires, err := usecase.issuer.issue(request, user)
	if err != nil {
		return nil, "", data.Expires{}, ErrIssueFailed
	}

	return ticket, nonce, expires, nil
}

type TokenIssuer struct {
	apiTokenIssuer     ApiTokenIssuer
	contentTokenIssuer ContentTokenIssuer
}

func NewTokenIssuer(api ApiTokenIssuer, content ContentTokenIssuer) TokenIssuer {
	return TokenIssuer{
		apiTokenIssuer:     api,
		contentTokenIssuer: content,
	}
}

func (usecase TokenIssuer) Issue(request data.Request, user data.User, expires data.Expires) (ApiToken, ContentToken, error) {
	apiToken, err := usecase.apiTokenIssuer.issue(request, user, expires)
	if err != nil {
		return nil, nil, ErrApiTokenIssueFailed
	}

	contentToken, err := usecase.contentTokenIssuer.issue(request, user, expires)
	if err != nil {
		return nil, nil, ErrContentTokenIssueFailed
	}

	return apiToken, contentToken, nil
}
