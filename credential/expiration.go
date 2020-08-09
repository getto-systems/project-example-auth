package credential

import (
	gotime "time"

	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/request"
)

type (
	ExtendLimit gotime.Time
	Expires     gotime.Time

	Expiration struct {
		expires     time.Second
		extendLimit time.Second
	}

	ExpirationParam struct {
		Expires     time.Second
		ExtendLimit time.Second
	}
)

func NewExpiration(param ExpirationParam) Expiration {
	return Expiration{
		expires:     param.Expires,
		extendLimit: param.ExtendLimit,
	}
}
func (exp Expiration) ExpireSecond() time.Second {
	return exp.expires
}
func (exp Expiration) Expires(request request.Request) Expires {
	return NewExpires(request.RequestedAt(), exp.expires)
}
func (exp Expiration) ExtendLimit(request request.Request) ExtendLimit {
	return NewExtendLimit(request.RequestedAt(), exp.extendLimit)
}

func NewExpires(requestedAt request.RequestedAt, second time.Second) Expires {
	return Expires(addSecond(requestedAt, second))
}
func NewExtendLimit(requestedAt request.RequestedAt, second time.Second) ExtendLimit {
	return ExtendLimit(addSecond(requestedAt, second))
}
func addSecond(requestedAt request.RequestedAt, second time.Second) gotime.Time {
	duration := gotime.Duration(second * 1_000_000_000)
	return gotime.Time(requestedAt).Add(duration)
}

func (expires Expires) Expired(requestedAt request.RequestedAt) bool {
	return gotime.Time(expires).Before(gotime.Time(requestedAt))
}

func (expires Expires) Limit(limit ExtendLimit) Expires {
	if gotime.Time(limit).Before(gotime.Time(expires)) {
		return Expires(limit)
	}
	return expires
}

func EmptyExpires() (_ Expires) {
	return
}
func EmptyExtendLimit() (_ ExtendLimit) {
	return
}
