package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type Verifier struct {
	pub  verifyEventPublisher
	repo verifyRepository
}

type verifyEventPublisher interface {
	VerifyPassword(data.Request, data.User)
	VerifyPasswordFailed(data.Request, data.User, error)
	AuthenticatedByPassword(data.Request, data.User)
}

type verifyDB interface {
	FindUserPassword(data.User) (data.HashedPassword, error)
}

func NewVerifier(
	pub verifyEventPublisher,
	db verifyDB,
	matcher Matcher,
) Verifier {
	return Verifier{
		pub:  pub,
		repo: newVerifyRepository(db, matcher),
	}
}

func (verifier Verifier) verify(request data.Request, user data.User, password data.RawPassword) error {
	verifier.pub.VerifyPassword(request, user)

	err := verifier.repo.matchPassword(user, password)
	if err != nil {
		verifier.pub.VerifyPasswordFailed(request, user, err)
		return err
	}

	verifier.pub.AuthenticatedByPassword(request, user)

	return nil
}

type verifyRepository struct {
	db      verifyDB
	matcher Matcher
}

func newVerifyRepository(db verifyDB, matcher Matcher) verifyRepository {
	return verifyRepository{
		db:      db,
		matcher: matcher,
	}
}

type Matcher interface {
	MatchPassword(data.HashedPassword, data.RawPassword) error
}

func (repo verifyRepository) matchPassword(user data.User, password data.RawPassword) error {
	hashed, err := repo.db.FindUserPassword(user)
	if err != nil {
		return err
	}

	return repo.matcher.MatchPassword(hashed, password)
}
