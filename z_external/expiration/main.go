package expiration

import (
	"time"
)

type (
	ExpireSecond int64
	Expires      time.Time

	ExtendSecond int64
	ExtendLimit  time.Time
)

func EmptyExpires() (_ Expires) {
	return
}
func EmptyExtendLimit() (_ ExtendLimit) {
	return
}

func NewExpires(now time.Time, second ExpireSecond) Expires {
	return Expires(addSecond(now, int64(second)))
}
func NewExtendLimit(now time.Time, second ExtendSecond) ExtendLimit {
	return ExtendLimit(addSecond(now, int64(second)))
}
func addSecond(now time.Time, second int64) time.Time {
	duration := time.Duration(second * 1_000_000_000)
	return time.Time(now).Add(duration)
}

func (limit ExtendLimit) Extend(now time.Time, second ExpireSecond) Expires {
	expires := NewExpires(now, second)

	if time.Time(limit).Before(time.Time(expires)) {
		expires = Expires(limit)
	}

	return expires
}

func (expires Expires) Expired(requestedAt time.Time) bool {
	return time.Time(expires).Before(requestedAt)
}
