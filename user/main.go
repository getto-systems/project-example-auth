package user

import (
	"strings"

	"github.com/getto-systems/project-example-id/basic"

	"errors"
	"fmt"
)

var (
	expireDuration = basic.Second(30)
	renewThreshold = basic.Second(5)

	ErrTicketHasNotEnoughPermission = errors.New("ticket has not enough permission")
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

func (user User) NewTicket(path basic.Path, requestedAt basic.RequestedAt) (basic.Ticket, error) {
	userID := user.userID
	roles := user.db.UserRoles(userID)

	authorized := requestedAt
	expires := requestedAt.Add(expireDuration)

	ticket := basic.Ticket{
		UserID:     userID,
		Roles:      roles,
		Authorized: authorized,
		Expires:    expires,
	}

	err := NewTicketInfo(ticket).HasEnoughPermission(path)
	if err != nil {
		return basic.Ticket{}, err
	}

	return ticket, nil
}

type TicketInfo struct {
	basic.Ticket
}

func NewTicketInfo(ticket basic.Ticket) TicketInfo {
	return TicketInfo{ticket}
}

func (ticket TicketInfo) IsRenewRequired(requestedAt basic.RequestedAt) bool {
	return ticket.Expires.Before(requestedAt.Add(renewThreshold))
}

func (ticket TicketInfo) HasEnoughPermission(path basic.Path) error {
	for _, role := range ticket.Roles {
		if role == super_role {
			return nil
		}

		if strings.HasPrefix(fmt.Sprintf("/%s/", role), string(path)) {
			return nil
		}
	}

	return ErrTicketHasNotEnoughPermission
}
