package user

import (
	"github.com/getto-systems/project-example-id/basic"

	"fmt"
)

var (
	expireSeconds  = basic.Second(30)
	renewThreshold = basic.Second(5)

	accessibleMap = AccessibleMap{}
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

func (f UserFactory) NewUser(userID basic.UserID) User {
	return User{f.db, userID}
}

type Ticket struct {
	userID     basic.UserID
	roles      basic.Roles
	authorized basic.Time
	expires    basic.Time
}

func (ticket Ticket) UserID() basic.UserID {
	return ticket.userID
}

func (ticket Ticket) Roles() basic.Roles {
	return ticket.roles
}

func (ticket Ticket) Authorized() basic.Time {
	return ticket.authorized
}

func (ticket Ticket) Expires() basic.Time {
	return ticket.expires
}

func (ticket Ticket) IsRenewRequired(now basic.Time) bool {
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

func (user User) NewTicket(path basic.Path, now basic.Time) (Ticket, error) {
	userID := user.userID
	roles := user.db.UserRoles(userID)

	authorized := now
	expires := now.Add(expireSeconds)

	return RestrictTicket(path, TicketData{
		UserID:     userID,
		Roles:      roles,
		Authorized: authorized,
		Expires:    expires,
	})
}

func RestrictTicket(path basic.Path, data TicketData) (Ticket, error) {
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

type TicketData struct {
	UserID     basic.UserID
	Roles      basic.Roles
	Authorized basic.Time
	Expires    basic.Time
}

func isAccessible(roles basic.Roles, path basic.Path) bool {
	for _, role := range roles {
		if role == super_role {
			return true
		}

		pathList, ok := accessibleMap[role]
		if ok {
			if pathList.contains(path) {
				return true
			}
		}
	}

	return false
}

type AccessibleMap map[string]PathList

type PathList []basic.Path

func (pathList PathList) contains(path basic.Path) bool {
	for _, accessible := range pathList {
		if path == accessible {
			return true
		}
	}

	return false
}
