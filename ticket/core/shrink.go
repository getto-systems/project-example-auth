package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"

	"errors"
)

type shrinker struct {
	pub  ticket.ShrinkEventPublisher
	repo shrinkRepository
}

func newShrinker(
	pub ticket.ShrinkEventPublisher,
	db ticket.ShrinkDB,
) shrinker {
	return shrinker{
		pub:  pub,
		repo: newShrinkRepository(db),
	}
}

func (shrinker shrinker) shrink(request data.Request, nonce ticket.Nonce, user data.User) error {
	shrinker.pub.ShrinkTicket(request, nonce, user)

	err := shrinker.repo.shrink(nonce, user)
	if err != nil {
		shrinker.pub.ShrinkTicketFailed(request, nonce, user, err)
		return err
	}

	return nil
}

type shrinkRepository struct {
	db ticket.ShrinkDB
}

func newShrinkRepository(db ticket.ShrinkDB) shrinkRepository {
	return shrinkRepository{
		db: db,
	}
}

func (repo shrinkRepository) shrink(nonce ticket.Nonce, user data.User) error {
	if !repo.db.TicketExists(nonce, user) {
		return errors.New("ticket not exists")
	}

	return repo.db.ShrinkTicket(nonce)
}
