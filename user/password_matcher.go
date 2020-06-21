package user

import (
	"github.com/getto-systems/project-example-id/data"
)

type PasswordMatcher interface {
	Match(password data.RawPassword) error
}

type userPasswordMatcher struct {
	engine PasswordMatchEngine
	hashed data.HashedPassword
}

type PasswordMatchEngine interface {
	MatchPassword(data.HashedPassword, data.RawPassword) error
}

func (m userPasswordMatcher) Match(password data.RawPassword) error {
	return m.engine.MatchPassword(m.hashed, password)
}

type notFoundPasswordMatcher struct {
	err error
}

func (m notFoundPasswordMatcher) Match(password data.RawPassword) error {
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

func (f PasswordMatcherFactory) New(hashed data.HashedPassword) PasswordMatcher {
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
