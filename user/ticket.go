package user

import (
	"github.com/getto-systems/project-example-id/basic"

	"errors"
)

var (
	expireDuration = basic.Second(30)
)

var (
	ErrTicketAlreadyExpired = errors.New("ticket already expired")
)

type Issuer struct {
	serializer TicketSerializer
	profile    basic.Profile
	original   basic.Ticket
}

type TicketSerializer interface {
	Serialize(basic.Ticket) (basic.Token, error)
}

func (iss Issuer) Authenticated(requestedAt basic.RequestedAt) (basic.Token, error) {
	expires := requestedAt.Expires(expireDuration)
	authenticatedAt := basic.AuthenticatedAt(requestedAt)

	token, err := iss.serializer.Serialize(basic.Ticket{
		Profile:         iss.profile,
		AuthenticatedAt: authenticatedAt,
		Expires:         expires,
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (iss Issuer) Renew(requestedAt basic.RequestedAt) (basic.Token, error) {
	if iss.original.Expires.Expired(requestedAt) {
		return nil, ErrTicketAlreadyExpired
	}

	expires := requestedAt.Expires(expireDuration)

	return iss.serializer.Serialize(basic.Ticket{
		Profile:         iss.profile,
		AuthenticatedAt: iss.original.AuthenticatedAt,
		Expires:         expires,
	})
}

type IssuerFactory struct {
	serializer TicketSerializer
}

func (f IssuerFactory) New(profile basic.Profile) Issuer {
	return Issuer{
		serializer: f.serializer,
		profile:    profile,
	}
}

func (f IssuerFactory) FromTicket(ticketInfo basic.Ticket) Issuer {
	return Issuer{
		serializer: f.serializer,
		profile:    ticketInfo.Profile,
		original:   ticketInfo,
	}
}

func NewIssuerFactory(serializer TicketSerializer) IssuerFactory {
	return IssuerFactory{
		serializer: serializer,
	}
}

type TicketAuthorizer struct {
	decoder TokenDecoder
	checker TicketPolicyChecker
}

type TokenDecoder interface {
	DecodeToken(basic.Token) (basic.Ticket, error)
}
type TicketPolicyChecker interface {
	HasEnoughPermission(basic.Ticket, basic.Request, basic.Resource) error
}

func NewTicketAuthorizer(decoder TokenDecoder, checker TicketPolicyChecker) TicketAuthorizer {
	return TicketAuthorizer{
		decoder: decoder,
		checker: checker,
	}
}

func (authorizer TicketAuthorizer) DecodeToken(token basic.Token) (basic.Ticket, error) {
	return authorizer.decoder.DecodeToken(token)
}

func (authorizer TicketAuthorizer) HasEnoughPermission(ticket basic.Ticket, request basic.Request, resource basic.Resource) error {
	return authorizer.checker.HasEnoughPermission(ticket, request, resource)
}
