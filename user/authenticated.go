package user

import (
	"github.com/getto-systems/project-example-id/data"
)

var (
	ticketExpireDuration = data.Second(30)
)

type UserAuthenticated struct {
	pub  UserAuthenticatedEventPublisher
	repo UserProfileRepository
	sign TicketSign

	userID  data.UserID
	request data.Request
}

func (user UserAuthenticated) IssueTicket() (data.Ticket, data.SignedTicket, error) {
	ticket := user.repo.Ticket(user.request, user.userID)

	signedTicket, err := user.sign.Sign(ticket)
	if err != nil {
		user.pub.TicketIssueFailed(user.request, ticket, err)
		return data.Ticket{}, nil, err
	}

	user.pub.Authenticated(user.request, ticket)

	return ticket, signedTicket, nil
}

type UserAuthenticatedEventPublisher interface {
	Authenticated(data.Request, data.Ticket)
	TicketIssueFailed(data.Request, data.Ticket, error)
}

type UserAuthenticatedEventHandler interface {
	UserAuthenticatedEventPublisher
}

type UserProfileRepository struct {
	db     UserProfileDB
	policy UserPermissionPolicy
}

func (repo UserProfileRepository) Ticket(request data.Request, userID data.UserID) data.Ticket {
	profile, err := repo.db.FindUserProfile(userID)
	if err != nil {
		// no permission when profile not found
		profile = data.Profile{
			UserID: userID,
		}
	}

	profile.Roles = repo.policy.Limit(request, profile.Roles)

	return NewTicket(request, profile)
}

func NewTicket(request data.Request, profile data.Profile) data.Ticket {
	requestedAt := request.RequestedAt
	expires := requestedAt.Expires(ticketExpireDuration)
	authenticatedAt := data.AuthenticatedAt(requestedAt)

	return data.Ticket{
		Profile:         profile,
		AuthenticatedAt: authenticatedAt,
		Expires:         expires,
	}
}

type UserProfileDB interface {
	FindUserProfile(data.UserID) (data.Profile, error)
}

type UserPermissionPolicy interface {
	Limit(data.Request, data.Roles) data.Roles
}

type UserAuthenticatedFactory struct {
	pub  UserAuthenticatedEventPublisher
	repo UserProfileRepository
	sign TicketSign
}

func NewUserAuthenticatedFactory(pub UserAuthenticatedEventPublisher, db UserProfileDB, policy UserPermissionPolicy, sign TicketSign) UserAuthenticatedFactory {
	return UserAuthenticatedFactory{
		pub: pub,
		repo: UserProfileRepository{
			db:     db,
			policy: policy,
		},
		sign: sign,
	}
}

func (f UserAuthenticatedFactory) New(request data.Request, userID data.UserID) UserAuthenticated {
	return UserAuthenticated{
		pub:  f.pub,
		repo: f.repo,
		sign: f.sign,

		userID:  userID,
		request: request,
	}
}
