package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

var (
	ErrVerifyFailed            = errors.New("ticket-verify-failed")
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
	IssueApiToken(data.Request, data.User, data.Roles, data.Expires)
	IssueApiTokenFailed(data.Request, data.User, data.Roles, data.Expires, error)

	IssueContentToken(data.Request, data.User, data.Expires)
	IssueContentTokenFailed(data.Request, data.User, data.Expires, error)

	ExtendTicket(data.Request, Nonce, data.User, data.Expires)
	ExtendTicketFailed(data.Request, Nonce, data.User, data.Expires, error)

	IssueTicket(data.Request, data.User, data.Expires, data.ExtendLimit)
	IssueTicketFailed(data.Request, data.User, data.Expires, data.ExtendLimit, error)

	ShrinkTicket(data.Request, Nonce, data.User)
	ShrinkTicketFailed(data.Request, Nonce, data.User, error)

	VerifyTicket(data.Request)
	VerifyTicketFailed(data.Request, error)
	AuthenticatedByTicket(data.Request, data.User)
}

type EventHandler interface {
	EventPublisher
}

type DB interface {
	FindUserRoles(data.User) (data.Roles, error)

	FindTicketExtendLimit(Nonce, data.User) (data.ExtendLimit, error)

	RegisterTransaction(Nonce, func(Nonce) error) (Nonce, error)
	RegisterTicket(Nonce, data.User, data.Expires, data.ExtendLimit) error
	NonceExists(Nonce) bool

	TicketExists(Nonce, data.User) bool
	ShrinkTicket(Nonce) error
}

type TicketExtender struct {
	verifier    Verifier
	extender    Extender
	tokenIssuer TokenIssuer
}

func NewTicketExtender(
	verifier Verifier,
	extender Extender,
	api ApiTokenIssuer,
	content ContentTokenIssuer,
) TicketExtender {
	return TicketExtender{
		verifier:    verifier,
		extender:    extender,
		tokenIssuer: NewTokenIssuer(api, content),
	}
}

func (usecase TicketExtender) Extend(request data.Request, ticket Ticket, nonce Nonce) (Ticket, ApiToken, ContentToken, data.Expires, error) {
	user, err := usecase.verifier.verify(request, ticket, nonce)
	if err != nil {
		return nil, nil, nil, data.Expires{}, ErrVerifyFailed
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

type TicketVerifier struct {
	verifier Verifier
}

func NewTicketVerifier(
	verifier Verifier,
) TicketVerifier {
	return TicketVerifier{
		verifier: verifier,
	}
}

func (usecase TicketVerifier) Verify(request data.Request, ticket Ticket, nonce Nonce) (data.User, error) {
	user, err := usecase.verifier.verify(request, ticket, nonce)
	if err != nil {
		return data.User{}, ErrVerifyFailed
	}
	return user, nil
}

type TicketShrinker struct {
	verifier Verifier
	shrinker Shrinker
}

func NewTicketShrinker(
	verifier Verifier,
	shrinker Shrinker,
) TicketShrinker {
	return TicketShrinker{
		verifier: verifier,
		shrinker: shrinker,
	}
}

func (usecase TicketShrinker) Shrink(request data.Request, ticket Ticket, nonce Nonce) error {
	user, err := usecase.verifier.verify(request, ticket, nonce)
	if err != nil {
		return ErrVerifyFailed
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
