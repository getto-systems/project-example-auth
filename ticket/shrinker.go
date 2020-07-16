package ticket

import (
	"errors"

	"github.com/getto-systems/project-example-id/data"
)

type Shrinker struct {
	pub  EventPublisher
	repo shrinkRepository
}

func NewShrinker(
	pub EventPublisher,
	db DB,
) Shrinker {
	return Shrinker{
		pub:  pub,
		repo: newShrinkRepository(db),
	}
}

func (shrinker Shrinker) shrink(request data.Request, nonce Nonce, user data.User) error {
	shrinker.pub.ShrinkTicket(request, nonce, user)

	err := shrinker.repo.shrink(nonce, user)
	if err != nil {
		shrinker.pub.ShrinkTicketFailed(request, nonce, user, err)
		return err
	}

	return nil
}

type shrinkRepository struct {
	db DB
}

func newShrinkRepository(db DB) shrinkRepository {
	return shrinkRepository{
		db: db,
	}
}

func (repo shrinkRepository) shrink(nonce Nonce, user data.User) error {
	if !repo.db.TicketExists(nonce, user) {
		return errors.New("ticket not exists")
	}

	return repo.db.ShrinkTicket(nonce)
}
