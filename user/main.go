package user

import (
	"fmt"
	"time"
)

var expireDuration = time.Duration(30 * 1_000_000_000)
var renewThreshold = time.Duration(5 * 1_000_000_000)

type User struct {
	db UserRepository

	userID UserID
}

type UserRepository interface {
	UserRoles(UserID) Roles
}

type UserID string
type Roles []string
type Path string

func (user User) UserID() UserID {
	return user.userID
}

type UserFactory struct {
	db UserRepository
}

func NewUserFactory(db UserRepository) UserFactory {
	return UserFactory{db}
}

func (f UserFactory) NewUser(userID UserID) User {
	return User{f.db, userID}
}

type Ticket struct {
	userID     UserID
	roles      Roles
	authorized time.Time
	expires    time.Time
}

func (ticket Ticket) UserID() UserID {
	return ticket.userID
}

func (ticket Ticket) Roles() Roles {
	return ticket.roles
}

func (ticket Ticket) Authorized() time.Time {
	return ticket.authorized
}

func (ticket Ticket) Expires() time.Time {
	return ticket.expires
}

func (ticket Ticket) IsRenewRequired(now time.Time) bool {
	return now.Before(ticket.Expires().Add(renewThreshold))
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

func (user User) NewTicket(path Path, now time.Time) (Ticket, error) {
	userID := user.userID
	roles := user.db.UserRoles(userID)

	authorized := now
	expires := now.Add(expireDuration)

	return RestrictTicket(path, TicketData{
		userID,
		roles,
		authorized,
		expires,
	})
}

func RestrictTicket(path Path, data TicketData) (Ticket, error) {
	var nullTicket Ticket

	if !data.Roles.isAccessible(path) {
		return nullTicket, fmt.Errorf("%s is not accessible as role: %v", path, data.Roles)
	}

	return Ticket{
		data.UserID,
		data.Roles,
		data.Authorized,
		data.Expires,
	}, nil
}

func (roles Roles) isAccessible(path Path) bool {
	return true // TODO fix: accessible check
}

type TicketData struct {
	UserID     UserID
	Roles      Roles
	Authorized time.Time
	Expires    time.Time
}
