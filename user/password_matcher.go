package user

import (
	"github.com/getto-systems/project-example-id/basic"
)

type PasswordMatcher interface {
	Match(password basic.RawPassword) error
}

type userPasswordMatcher struct {
	engine PasswordMatchEngine
	hashed basic.HashedPassword
}

type PasswordMatchEngine interface {
	MatchPassword(basic.HashedPassword, basic.RawPassword) error
}

func (m userPasswordMatcher) Match(password basic.RawPassword) error {
	return m.engine.MatchPassword(m.hashed, password)
}

type notFoundPasswordMatcher struct {
	err error
}

func (m notFoundPasswordMatcher) Match(password basic.RawPassword) error {
	return m.err
}

type PasswordMatcherFactory struct {
	engine PasswordMatchEngine
}

func NewPasswordMatcherFactory(engine PasswordMatchEngine) PasswordMatcherFactory {
	return PasswordMatcherFactory{
		engine: engine,
	}
}

func (f PasswordMatcherFactory) New(hashed basic.HashedPassword) PasswordMatcher {
	return userPasswordMatcher{
		engine: f.engine,
		hashed: hashed,
	}
}

func (f PasswordMatcherFactory) NotFound(err error) PasswordMatcher {
	return notFoundPasswordMatcher{
		err: err,
	}
}
