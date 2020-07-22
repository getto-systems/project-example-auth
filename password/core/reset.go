package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
)

type resetter struct {
	logger password.ResetLogger
	exp    password.Expiration
	repo   resetRepository
}

func newResetter(
	logger password.ResetLogger,
	db password.ResetDB,
	gen password.ResetGenerator,
) resetter {
	return resetter{
		logger: logger,
		repo:   newResetRepository(db, gen),
	}
}

func (resetter resetter) issueReset(request data.Request, login password.Login) (password.Reset, error) {
	expires := resetter.exp.Expires(request)
	resetter.logger.TryToIssueReset(request, login, expires)

	// TODO 第2引数の token は notifier に渡す
	reset, _, user, err := resetter.repo.register(login, request, expires)
	if err != nil {
		resetter.logger.FailedToIssueReset(request, login, expires, err)
		return password.Reset{}, err
	}

	resetter.logger.IssuedReset(request, reset, user, expires)

	return reset, nil
}

func (resetter resetter) getResetStatus(request data.Request, reset password.Reset) (password.ResetStatus, error) {
	resetter.logger.TryToGetResetStatus(request, reset)

	status, err := resetter.repo.findResetStatus(reset)
	if err != nil {
		resetter.logger.FailedToGetResetStatus(request, reset, err)
		return password.ResetStatus{}, err
	}

	return status, nil
}

func (resetter resetter) validate(request data.Request, login password.Login, token password.ResetToken) (data.User, error) {
	resetter.logger.TryToValidateResetToken(request)

	user, err := resetter.repo.findUserByResetToken(request, login, token)
	if err != nil {
		resetter.logger.FailedToValidateResetToken(request, err)
		return data.User{}, err
	}

	resetter.logger.AuthedByResetToken(request, user)

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

func (repo resetRepository) register(login password.Login, request data.Request, expires data.Expires) (password.Reset, password.ResetToken, data.User, error) {
	user, err := repo.findUserByLogin(login)
	if err != nil {
		return password.Reset{}, "", data.User{}, err
	}

	reset, token, err := repo.db.RegisterReset(repo.gen, user, request.RequestedAt(), expires)
	if err != nil {
		return password.Reset{}, "", data.User{}, err
	}

	return reset, token, user, nil
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
