package basic

import (
	"fmt"
	"time"
)

type (
	UserID string
	Roles  []string

	Path string

	RequestedAt time.Time
	Expires     time.Time
	Second      int64

	Ticket struct {
		UserID     UserID
		Roles      Roles
		Authorized RequestedAt
		Expires    Expires
	}
)

func (requestedAt RequestedAt) Add(seconds Second) Expires {
	duration := time.Duration(seconds * 1_000_000_000)
	after := time.Time(requestedAt).Add(duration)
	return Expires(after)
}

func (requestedAt RequestedAt) String() string {
	return time.Time(requestedAt).String()
}

func (expires Expires) Before(target Expires) bool {
	return time.Time(expires).Before(time.Time(target))
}

func (expires Expires) String() string {
	return time.Time(expires).String()
}

func Mimute(minutes int64) Second {
	return Second(minutes * 60)
}

func (ticket Ticket) String() string {
	return fmt.Sprintf(
		"Ticket{UserID:%s, Roles:%s, Authorized:%s, Expires:%s}",
		ticket.UserID,
		ticket.Roles,
		ticket.Authorized.String(),
		ticket.Expires.String(),
	)
}
