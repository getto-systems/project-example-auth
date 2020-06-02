package basic

import (
	"time"
)

type (
	UserID string
	Roles  []string

	Path string

	Time   time.Time
	Second int64
)

func (t Time) Add(seconds Second) Time {
	duration := time.Duration(seconds * 1_000_000_000)
	after := time.Time(t).Add(duration)
	return Time(after)
}

func (t Time) Before(target Time) bool {
	return time.Time(t).Before(time.Time(target))
}

func (t Time) String() string {
	return time.Time(t).String()
}
