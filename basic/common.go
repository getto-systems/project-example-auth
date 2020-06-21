package basic

import (
	"fmt"
	"time"
)

type (
	UserID string
	Roles  []string

	Profile struct {
		UserID UserID
		Roles  Roles
	}

	RequestedAt     time.Time
	Expires         time.Time
	AuthenticatedAt time.Time
	Second          int64

	Path       string
	RemoteAddr string

	Resource struct {
		Path Path
	}
	Route struct {
		RemoteAddr RemoteAddr
	}
	Request struct {
		RequestedAt RequestedAt
		Route       Route
	}

	Token  []byte
	Ticket struct {
		Profile         Profile
		AuthenticatedAt AuthenticatedAt
		Expires         Expires
	}
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
	return time.Time(now).Before(time.Time(expires))
}

func (requestedAt RequestedAt) String() string {
	return requestedAt.Time().String()
}

func (requestedAt RequestedAt) Time() time.Time {
	return time.Time(requestedAt)
}

func (expires Expires) String() string {
	return expires.Time().String()
}

func (expires Expires) Time() time.Time {
	return time.Time(expires)
}

func (authenticatedAt AuthenticatedAt) String() string {
	return authenticatedAt.Time().String()
}

func (authenticatedAt AuthenticatedAt) Time() time.Time {
	return time.Time(authenticatedAt)
}

func (profile Profile) String() string {
	return fmt.Sprintf(
		"Profile{UserID:%s, Roles:%s}",
		profile.UserID,
		profile.Roles,
	)
}

func (info Ticket) String() string {
	return fmt.Sprintf(
		"Ticket{Profile:%s, Authorized:%s, Expires:%s}",
		info.Profile,
		info.AuthenticatedAt.String(),
		info.Expires.String(),
	)
}
