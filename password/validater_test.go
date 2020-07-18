package password

import (
	"errors"
	"strings"
	"testing"

	"github.com/getto-systems/project-example-id/data"
)

// パスワードが一致したら AuthenticatedByPassword イベントが発行される
func TestValidate(t *testing.T) {
	hashed := data.HashedPassword("password")
	raw := data.RawPassword("password")

	pub := newValidateTestEventPublisher()
	db := newValidateTestDB(hashed)
	matcher := newValidateTestMatcher()

	// validate
	validater := NewValidater(pub, db, matcher)
	err := validater.validate(data.Request{}, data.User{}, raw)

	h := newValidateTestHelper(t, pub, err)
	h.checkValidateError(nil)
	h.checkValidatePasswordEvent("fired")
	h.checkValidatePasswordFailedEvent(nil)
	h.checkAuthenticatedByPasswordEvent("fired")
}

// パスワードが一致しなかったら ValidatePasswordFailed イベントが発行される
func TestValidateFailed(t *testing.T) {
	hashed := data.HashedPassword("password")
	raw := data.RawPassword("different-password")

	pub := newValidateTestEventPublisher()
	db := newValidateTestDB(hashed)
	matcher := newValidateTestMatcher()

	// validate
	validater := NewValidater(pub, db, matcher)
	err := validater.validate(data.Request{}, data.User{}, raw)

	h := newValidateTestHelper(t, pub, err)
	h.checkValidateError(errors.New("password not matched"))
	h.checkValidatePasswordEvent("fired")
	h.checkValidatePasswordFailedEvent(errors.New("password not matched"))
	h.checkAuthenticatedByPasswordEvent("never")
}

// パスワードが見つからない場合、必ず失敗する
func TestValidateFailedWhenPasswordNotFound(t *testing.T) {
	raw := data.RawPassword("password")

	pub := newValidateTestEventPublisher()
	db := emptyValidateTestDB()
	matcher := newValidateTestMatcher()

	// validate
	validater := NewValidater(pub, db, matcher)
	err := validater.validate(data.Request{}, data.User{}, raw)

	h := newValidateTestHelper(t, pub, err)
	h.checkValidateError(errors.New("password not found"))
	h.checkValidatePasswordEvent("fired")
	h.checkValidatePasswordFailedEvent(errors.New("password not found"))
	h.checkAuthenticatedByPasswordEvent("never")
}

// 空のパスワードの場合、必ず失敗する
func TestValidateFailedWhenEmptyPassword(t *testing.T) {
	hashed := data.HashedPassword("")
	raw := data.RawPassword("")

	pub := newValidateTestEventPublisher()
	db := newValidateTestDB(hashed)
	matcher := newValidateTestMatcher()

	// validate
	validater := NewValidater(pub, db, matcher)
	err := validater.validate(data.Request{}, data.User{}, raw)

	h := newValidateTestHelper(t, pub, err)
	h.checkValidateError(errors.New("password is empty"))
	h.checkValidatePasswordEvent("fired")
	h.checkValidatePasswordFailedEvent(errors.New("password is empty"))
	h.checkAuthenticatedByPasswordEvent("never")
}

// 長いパスワードの場合、必ず失敗する
func TestValidateFailedWhenLongPassword(t *testing.T) {
	length := 73
	hashed := data.HashedPassword(strings.Repeat("a", length))
	raw := data.RawPassword(strings.Repeat("a", length))

	pub := newValidateTestEventPublisher()
	db := newValidateTestDB(hashed)
	matcher := newValidateTestMatcher()

	// validate
	validater := NewValidater(pub, db, matcher)
	err := validater.validate(data.Request{}, data.User{}, raw)

	h := newValidateTestHelper(t, pub, err)
	h.checkValidateError(errors.New("password is too long"))
	h.checkValidatePasswordEvent("fired")
	h.checkValidatePasswordFailedEvent(errors.New("password is too long"))
	h.checkAuthenticatedByPasswordEvent("never")
}

// ギリギリの長さのパスワードの場合、成功する
func TestValidateWhenLongPassword(t *testing.T) {
	length := 72
	hashed := data.HashedPassword(strings.Repeat("a", length))
	raw := data.RawPassword(strings.Repeat("a", length))

	pub := newValidateTestEventPublisher()
	db := newValidateTestDB(hashed)
	matcher := newValidateTestMatcher()

	// validate
	validater := NewValidater(pub, db, matcher)
	err := validater.validate(data.Request{}, data.User{}, raw)

	h := newValidateTestHelper(t, pub, err)
	h.checkValidateError(nil)
	h.checkValidatePasswordEvent("fired")
	h.checkValidatePasswordFailedEvent(nil)
	h.checkAuthenticatedByPasswordEvent("fired")
}

type (
	validateTestEventPublisher struct {
		validatePassword        string
		validatePasswordFailed  error
		authenticatedByPassword string
	}

	validateTestDB struct {
		password *data.HashedPassword
	}

	validateTestMatcher struct {
	}

	validateTestHelper struct {
		t   *testing.T
		pub *validateTestEventPublisher
		err error
	}
)

func newValidateTestEventPublisher() *validateTestEventPublisher {
	return &validateTestEventPublisher{
		validatePassword:        "never",
		authenticatedByPassword: "never",
	}
}

func (pub *validateTestEventPublisher) ValidatePassword(request data.Request, user data.User) {
	pub.validatePassword = "fired"
}
func (pub *validateTestEventPublisher) ValidatePasswordFailed(request data.Request, user data.User, err error) {
	pub.validatePasswordFailed = err
}
func (pub *validateTestEventPublisher) AuthenticatedByPassword(request data.Request, user data.User) {
	pub.authenticatedByPassword = "fired"
}

func newValidateTestDB(password data.HashedPassword) validateTestDB {
	return validateTestDB{
		password: &password,
	}
}

func emptyValidateTestDB() validateTestDB {
	return validateTestDB{}
}

func (db validateTestDB) FindUserPassword(user data.User) (data.HashedPassword, error) {
	if db.password == nil {
		return nil, errors.New("password not found")
	}
	return *db.password, nil
}

func newValidateTestMatcher() validateTestMatcher {
	return validateTestMatcher{}
}

func (validateTestMatcher) MatchPassword(hashed data.HashedPassword, raw data.RawPassword) error {
	if string(raw) != string(hashed) {
		return errors.New("password not matched")
	}
	return nil
}

func newValidateTestHelper(t *testing.T, pub *validateTestEventPublisher, err error) validateTestHelper {
	return validateTestHelper{
		t:   t,
		pub: pub,
		err: err,
	}
}

func (h validateTestHelper) checkValidateError(err error) {
	if err == nil {
		if h.err != nil {
			h.t.Errorf("validate fired: %s", h.err)
		}
	} else {
		if h.err == nil {
			h.t.Error("validate success")
		} else {
			if h.err.Error() != err.Error() {
				h.t.Errorf("validate error message is not matched: %s (expected: %s)", h.err, err)
			}
		}
	}
}
func (h validateTestHelper) checkValidatePasswordEvent(event string) {
	if h.pub.validatePassword != event {
		h.t.Errorf("ValidatePassword event not match: %s (expected: %s)", h.pub.validatePassword, event)
	}
}
func (h validateTestHelper) checkValidatePasswordFailedEvent(err error) {
	if err == nil {
		if h.pub.validatePasswordFailed != nil {
			h.t.Errorf("ValidatePasswordFailed event fired: %s", h.pub.validatePasswordFailed)
		}
	} else {
		if h.pub.validatePasswordFailed == nil {
			h.t.Error("ValidatePasswordFailed event not fired")
		} else {
			if h.pub.validatePasswordFailed.Error() != err.Error() {
				h.t.Errorf("ValidatePasswordFailed error message is not matched: %s (expected: %s)", h.pub.validatePasswordFailed, err)
			}
		}
	}
}
func (h validateTestHelper) checkAuthenticatedByPasswordEvent(event string) {
	if h.pub.authenticatedByPassword != event {
		h.t.Errorf("AuthenticatedByPassword event not match: %s (expected: %s)", h.pub.authenticatedByPassword, event)
	}
}
