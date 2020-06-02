package basic

import (
	"time"
)

type (
	UserID string
	Roles  []string

	Path string

	RequestedAt time.Time
	Expires     time.Time
	Second      int64
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
