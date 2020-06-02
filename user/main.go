package user

import (
	"fmt"
	"time"
)

var (
	expireDuration = time.Duration(30 * 1_000_000_000)
	renewThreshold = time.Duration(5 * 1_000_000_000)

	accessibleMap = AccessibleMap{}
)

const (
	super_role = "admin"
)

type User struct {
	db UserRepository

	userID UserID
}

type UserRepository interface {
	UserRoles(UserID) Roles
}

type (
	UserID string
	Roles  []string
	Path   string
)

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
		UserID:     userID,
		Roles:      roles,
		Authorized: authorized,
		Expires:    expires,
	})
}

func RestrictTicket(path Path, data TicketData) (Ticket, error) {
	if !data.Roles.isAccessible(path) {
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
	UserID     UserID
	Roles      Roles
	Authorized time.Time
	Expires    time.Time
}

func (roles Roles) isAccessible(path Path) bool {
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

type PathList []Path

func (pathList PathList) contains(path Path) bool {
	for _, accessible := range pathList {
		if path == accessible {
			return true
		}
	}

	return false
}
