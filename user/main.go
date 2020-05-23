package user

import (
	"fmt"
	"time"
)

var renewThreshold = time.Duration(5 * 1_000_000_000)
var expireDuration = time.Duration(30 * 1_000_000_000)

type User struct {
	userID UserID
	roles  Roles
}

type UserRepository interface {
	UserRoles(UserID) Roles
}

type UserID string
type Roles []string

func (user User) UserID() UserID {
	return user.userID
}

func (user User) Roles() Roles {
	return user.roles
}

func (user User) NewTicket(now time.Time) Ticket {
	expires := now.Add(expireDuration)
	return Ticket{user, expires}
}

type Ticket struct {
	user    User
	expires time.Time
}

func (ticket Ticket) User() User {
	return ticket.user
}

func (ticket Ticket) Expires() time.Time {
	return ticket.expires
}

func (ticket Ticket) IsRenewRequired(now time.Time) bool {
	return now.Before(ticket.Expires().Add(renewThreshold))
}

func NewUser(userID UserID, roles Roles, path Path) (User, error) {
	if !roles.isAccessible(path) {
		return User{}, fmt.Errorf("%s is not accessible in role: [%v]", path, roles)
	}

	return User{userID, roles}, nil
}

type Path string

func (roles Roles) isAccessible(path Path) bool {
	return true // TODO fix: accessible check
}
