package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type resetter struct {
	pub  password.ResetEventPublisher
	exp  password.Expiration
	repo resetRepository
}

func newResetter(
	pub password.ResetEventPublisher,
	db password.ResetDB,
	gen password.ResetGenerator,
) resetter {
	return resetter{
		pub:  pub,
		repo: newResetRepository(db, gen),
	}
}

func (resetter resetter) issueResetToken(request data.Request, login password.Login) (password.Reset, error) {
	expires := resetter.exp.Expires(request)
	resetter.pub.IssueResetToken(request, login, expires)

	reset, err := resetter.repo.register(login, request, expires)
	if err != nil {
		resetter.pub.IssueResetTokenFailed(request, login, expires, err)
		return password.Reset{}, err
	}

	return reset, nil
}

func (resetter resetter) getResetStatus(request data.Request, reset password.Reset) (password.ResetStatus, error) {
	resetter.pub.GetResetStatus(request, reset)

	status, err := resetter.repo.findResetStatus(reset)
	if err != nil {
		resetter.pub.GetResetStatusFailed(request, reset, err)
		return password.ResetStatus{}, err
	}

	return status, nil
}

func (resetter resetter) validate(request data.Request, login password.Login, token password.ResetToken) (data.User, error) {
	resetter.pub.ValidateResetToken(request)

	user, err := resetter.repo.findUserByResetToken(request, login, token)
	if err != nil {
		resetter.pub.ValidateResetTokenFailed(request, err)
		return data.User{}, err
	}

	resetter.pub.AuthenticatedByResetToken(request, user)

	return user, nil
}

type resetRepository struct {
	db  password.ResetDB
	gen password.ResetGenerator
}

func newResetRepository(db password.ResetDB, gen password.ResetGenerator) resetRepository {
	return resetRepository{
		db:  db,
		gen: gen,
	}
}

func (repo resetRepository) register(login password.Login, request data.Request, expires data.Expires) (password.Reset, error) {
	user, err := repo.findUserByLogin(login)
	if err != nil {
		return password.Reset{}, err
	}

	return repo.db.RegisterReset(repo.gen, user, request.RequestedAt(), expires)
}

func (repo resetRepository) findUserByResetToken(request data.Request, login password.Login, token password.ResetToken) (data.User, error) {
	loginUser, err := repo.findUserByLogin(login)
	if err != nil {
		return data.User{}, err
	}

	user, err := repo.findResetUser(token)
	if err != nil {
		return data.User{}, err
	}

	err = user.Validate(request, loginUser)
	if err != nil {
		return data.User{}, err
	}

	return user.User(), nil
}

func (repo resetRepository) findUserByLogin(login password.Login) (data.User, error) {
	userSlice, err := repo.db.FilterUserByLogin(login)
	if err != nil {
		return data.User{}, err
	}

	if len(userSlice) == 0 {
		return data.User{}, password.ErrResetUserNotFound
	}

	return userSlice[0], nil
}

func (repo resetRepository) findResetStatus(reset password.Reset) (password.ResetStatus, error) {
	statusSlice, err := repo.db.FilterResetStatus(reset)
	if err != nil {
		return password.ResetStatus{}, err
	}

	if len(statusSlice) == 0 {
		// ステータスが見つからない場合は「リクエストなし」を返してエラーにはしない
		return password.NewResetStatusNotRequested(), nil
	}

	return statusSlice[0], nil
}

func (repo resetRepository) findResetUser(token password.ResetToken) (password.ResetUser, error) {
	userSlice, err := repo.db.FilterResetUser(token)
	if err != nil {
		return password.ResetUser{}, err
	}

	if len(userSlice) == 0 {
		return password.ResetUser{}, password.ErrResetTokenNotFound
	}

	return userSlice[0], nil
}
