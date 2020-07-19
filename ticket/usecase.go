package ticket

import (
	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	ErrValidateFailed          = errors.New("ticket-validate-failed")
	ErrExtendFailed            = errors.New("ticket-extend-failed")
	ErrShrinkFailed            = errors.New("ticket-shrink-failed")
	ErrIssueFailed             = errors.New("ticket-issue-failed")
	ErrApiTokenIssueFailed     = errors.New("api-token-issue-failed")
	ErrContentTokenIssueFailed = errors.New("content-token-issue-failed")
)

type EventPublisher interface {
	IssueEventPublisher
	ExtendEventPublisher
	ValidateEventPublisher
	ShrinkEventPublisher

	IssueApiTokenEventPublisher
	IssueContentTokenEventPublisher
}

type EventHandler interface {
	EventPublisher
}

type DB interface {
	IssueDB
	ExtendDB
	ShrinkDB

	IssueApiTokenDB
}

type Usecase interface {
	Issue(request data.Request, user data.User) (Ticket, Nonce, data.Expires, error)
	Extend(request data.Request, ticket Ticket, nonce Nonce) (Ticket, ApiToken, ContentToken, data.Expires, error)
	Validate(request data.Request, ticket Ticket, nonce Nonce) (data.User, error)
	Shrink(request data.Request, ticket Ticket, nonce Nonce) error

	IssueToken(request data.Request, user data.User, expires data.Expires) (ApiToken, ContentToken, error)
}
