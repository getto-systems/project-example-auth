package time

import (
	"time"
)

type (
	RequestedAt time.Time
	ExtendLimit time.Time
	Expires     time.Time
	Second      int64

	Time time.Time
)

func Now() RequestedAt {
	return RequestedAt(time.Now().UTC())
}

func Minute(minutes int64) Second {
	return Second(minutes * 60)
}
func Hour(hours int64) Second {
	return Minute(hours * 60)
}
func Day(days int64) Second {
	return Hour(days * 24)
}

func (requestedAt RequestedAt) Expires(second Second) Expires {
	return Expires(requestedAt.addSecond(second))
}
func (requestedAt RequestedAt) ExtendLimit(second Second) ExtendLimit {
	return ExtendLimit(requestedAt.addSecond(second))
}
func (requestedAt RequestedAt) addSecond(second Second) time.Time {
	duration := time.Duration(second * 1_000_000_000)
	return time.Time(requestedAt).Add(duration)
}

func (requestedAt RequestedAt) Time() Time {
	return Time(requestedAt)
}

func (requestedAt RequestedAt) Expired(expires Expires) bool {
	return time.Time(expires).Before(time.Time(requestedAt))
}

func (expires Expires) Limit(limit ExtendLimit) Expires {
	if time.Time(limit).Before(time.Time(expires)) {
		return Expires(limit)
	}
	return expires
}

func EmptyExtendLimit() (_ ExtendLimit) {
	return
}
