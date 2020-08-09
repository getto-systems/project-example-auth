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

func ExpireMinute(minutes int64) ExpireSecond {
	// 有効期限は短めにするべきなので minute まで提供
	return ExpireSecond(minutes * 60)
}

func ExtendMinute(minutes int64) ExtendSecond {
	return ExtendSecond(minutes * 60)
}
func ExtendHour(hours int64) ExtendSecond {
	return ExtendMinute(hours * 60)
}
func ExtendDay(days int64) ExtendSecond {
	return ExtendHour(days * 24)
}

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

func Extend(now time.Time, limit ExtendLimit, second ExpireSecond) Expires {
	expires := NewExpires(now, second)

	if time.Time(limit).Before(time.Time(expires)) {
		expires = Expires(limit)
	}

	return expires
}

func (expires Expires) Expired(requestedAt time.Time) bool {
	return time.Time(expires).Before(requestedAt)
}
