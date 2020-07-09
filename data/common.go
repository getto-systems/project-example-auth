package data

import (
	"time"
)

type (
	UserID string
	Roles  []string

	User struct {
		UserID UserID
	}

	Profile struct {
		UserID UserID
		Roles  Roles
	}

	RequestedAt     time.Time
	Expires         time.Time
	AuthenticatedAt time.Time
	Second          int64

	RemoteAddr string

	Route struct {
		RemoteAddr RemoteAddr
	}
	Request struct {
		RequestedAt RequestedAt
		Route       Route
	}

	Ticket struct {
		Profile         Profile
		AuthenticatedAt AuthenticatedAt
		Expires         Expires
	}

	SignedTicket []byte
)

func Mimute(minutes int64) Second {
	return Second(minutes * 60)
}

func (requestedAt RequestedAt) Expires(seconds Second) Expires {
	duration := time.Duration(seconds * 1_000_000_000)
	expires := time.Time(requestedAt).Add(duration)
	return Expires(expires)
}

func (expires Expires) Expired(now RequestedAt) bool {
	return time.Time(expires).Before(time.Time(now))
}

func (requestedAt RequestedAt) String() string {
	return time.Time(requestedAt).String()
}

func (expires Expires) String() string {
	return expires.Time().String()
}

func (expires Expires) Time() time.Time {
	return time.Time(expires)
}

func (authenticatedAt AuthenticatedAt) String() string {
	return time.Time(authenticatedAt).String()
}
