package user

import (
	"github.com/getto-systems/project-example-id/data"

	"errors"
)

var (
	expireDuration = data.Second(30)
)

var (
	ErrTicketAlreadyExpired = errors.New("ticket already expired")
)

type Issuer struct {
	serializer TicketSerializer
	profile    data.Profile
	original   data.Ticket
}

type TicketSerializer interface {
	Serialize(data.Ticket) (data.Token, error)
}

func (iss Issuer) Authenticated(requestedAt data.RequestedAt) (data.Token, error) {
	expires := requestedAt.Expires(expireDuration)
	authenticatedAt := data.AuthenticatedAt(requestedAt)

	token, err := iss.serializer.Serialize(data.Ticket{
		Profile:         iss.profile,
		AuthenticatedAt: authenticatedAt,
		Expires:         expires,
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (iss Issuer) Renew(requestedAt data.RequestedAt) (data.Token, error) {
	if iss.original.Expires.Expired(requestedAt) {
		return nil, ErrTicketAlreadyExpired
	}

	expires := requestedAt.Expires(expireDuration)

	return iss.serializer.Serialize(data.Ticket{
		Profile:         iss.profile,
		AuthenticatedAt: iss.original.AuthenticatedAt,
		Expires:         expires,
	})
}

type IssuerFactory struct {
	serializer TicketSerializer
}

func (f IssuerFactory) New(profile data.Profile) Issuer {
	return Issuer{
		serializer: f.serializer,
		profile:    profile,
	}
}

func (f IssuerFactory) FromTicket(ticketInfo data.Ticket) Issuer {
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
	DecodeToken(data.Token) (data.Ticket, error)
}
type TicketPolicyChecker interface {
	HasEnoughPermission(data.Ticket, data.Request, data.Resource) error
}

func NewTicketAuthorizer(decoder TokenDecoder, checker TicketPolicyChecker) TicketAuthorizer {
	return TicketAuthorizer{
		decoder: decoder,
		checker: checker,
	}
}

func (authorizer TicketAuthorizer) DecodeToken(token data.Token) (data.Ticket, error) {
	return authorizer.decoder.DecodeToken(token)
}

func (authorizer TicketAuthorizer) HasEnoughPermission(ticket data.Ticket, request data.Request, resource data.Resource) error {
	return authorizer.checker.HasEnoughPermission(ticket, request, resource)
}
