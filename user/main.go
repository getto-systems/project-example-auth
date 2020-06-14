package user

import (
	"strings"

	"github.com/getto-systems/project-example-id/basic"

	"fmt"
)

var (
	expireDuration = basic.Second(30)
	renewThreshold = basic.Second(5)
)

const (
	super_role = "admin"
)

type User struct {
	db UserRepository

	userID basic.UserID
}

type UserRepository interface {
	UserRoles(basic.UserID) basic.Roles
}

func (user User) UserID() basic.UserID {
	return user.userID
}

type UserFactory struct {
	db UserRepository
}

func NewUserFactory(db UserRepository) UserFactory {
	return UserFactory{db}
}

func (f UserFactory) New(userID basic.UserID) User {
	return User{f.db, userID}
}

type Ticket struct {
	userID     basic.UserID
	roles      basic.Roles
	authorized basic.RequestedAt
	expires    basic.Expires
}

func (ticket Ticket) IsRenewRequired(requestedAt basic.RequestedAt) bool {
	return ticket.Expires().Before(requestedAt.Add(renewThreshold))
}

func (ticket Ticket) UserID() basic.UserID {
	return ticket.userID
}

func (ticket Ticket) Expires() basic.Expires {
	return ticket.expires
}

func (ticket Ticket) Data() basic.TicketData {
	return basic.TicketData{
		UserID:     ticket.userID,
		Roles:      ticket.roles,
		Authorized: ticket.authorized,
		Expires:    ticket.expires,
	}
}

func (ticket Ticket) String() string {
	return fmt.Sprintf(
		"Ticket{userID:%s, roles:%s, authorized:%s, expires:%s}",
		ticket.userID,
		ticket.roles,
		ticket.authorized.String(),
		ticket.expires.String(),
	)
}

func (user User) NewTicket(path basic.Path, requestedAt basic.RequestedAt) (Ticket, error) {
	userID := user.userID
	roles := user.db.UserRoles(userID)

	authorized := requestedAt
	expires := requestedAt.Add(expireDuration)

	return RestrictTicket(path, basic.TicketData{
		UserID:     userID,
		Roles:      roles,
		Authorized: authorized,
		Expires:    expires,
	})
}

func RestrictTicket(path basic.Path, data basic.TicketData) (Ticket, error) {
	if !isAccessible(data.Roles, path) {
		return Ticket{}, fmt.Errorf("%s is not accessible as role: %v", path, data.Roles)
	}

	return Ticket{
		userID:     data.UserID,
		roles:      data.Roles,
		authorized: data.Authorized,
		expires:    data.Expires,
	}, nil
}

func isAccessible(roles basic.Roles, path basic.Path) bool {
	for _, role := range roles {
		if role == super_role {
			return true
		}

		if strings.HasPrefix(fmt.Sprintf("/%s/", role), string(path)) {
			return true
		}
	}

	return false
}
