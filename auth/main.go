package auth

import (
	"fmt"
	"time"
)

var renewThreshold = time.Duration(5 * 1_000_000_000)
var expireDuration = time.Duration(30 * 1_000_000_000)

type Tokener interface {
	Parse(TicketToken, Path) (*Ticket, error)
	TicketToken(*Ticket) (TicketToken, error)
	FullToken(*Ticket) ([]byte, error)
}

type TicketToken string

type Ticket struct {
	user    *User
	expires time.Time
}

func (ticket Ticket) User() *User {
	return ticket.user
}

func (ticket Ticket) Expires() time.Time {
	return ticket.expires
}

type User struct {
	userID UserID
	roles  Roles
}

func (user *User) UserID() UserID {
	return user.userID
}

func (user *User) Roles() Roles {
	return user.roles
}

func (ticket Ticket) IsRenewRequired(now time.Time) bool {
	return now.Before(ticket.Expires().Add(renewThreshold))
}

type UserID string
type Roles []string
type Path string

type UserRepository interface {
	UserRoles(UserID) Roles
}

func NewUserFromRepository(db UserRepository, userID UserID, path Path) (*User, error) {
	return NewUser(
		userID,
		db.UserRoles(userID),
		path,
	)
}

func NewUser(userID UserID, roles Roles, path Path) (*User, error) {
	if !roles.isAccessible(path) {
		return nil, fmt.Errorf("%s is not accessible in role: [%v]", path, roles)
	}

	return &User{userID, roles}, nil
}

func (roles Roles) isAccessible(path Path) bool {
	return true // TODO fix: accessible check
}

func (user *User) NewTicket(now time.Time) *Ticket {
	expires := now.Add(expireDuration)
	return &Ticket{user, expires}
}
