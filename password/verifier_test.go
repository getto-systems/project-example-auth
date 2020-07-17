package password

import (
	"errors"
	"testing"

	"github.com/getto-systems/project-example-id/data"
)

// パスワードが一致したら AuthenticatedByPassword イベントが発行される
func TestVerify(t *testing.T) {
	hashed := data.HashedPassword("password")
	raw := data.RawPassword("password")

	pub := newVerifyTestEventPublisher()
	db := newVerifyTestDB(hashed)
	matcher := newVerifyTestMatcher()

	// verify
	verifier := NewVerifier(pub, db, matcher)
	err := verifier.verify(data.Request{}, data.User{}, raw)

	h := newVerifyTestHelper(t, pub, err)
	h.checkVerifyError(nil)
	h.checkVerifyPasswordEvent("VerifyPassword")
	h.checkVerifyPasswordFailedEvent(nil)
	h.checkAuthenticatedByPasswordEvent("AuthenticatedByPassword")
}

// パスワードが一致しなかったら VerifyPasswordFailed イベントが発行される
func TestVerifyFailed(t *testing.T) {
	hashed := data.HashedPassword("password")
	raw := data.RawPassword("different-password")

	pub := newVerifyTestEventPublisher()
	db := newVerifyTestDB(hashed)
	matcher := newVerifyTestMatcher()

	// verify
	verifier := NewVerifier(pub, db, matcher)
	err := verifier.verify(data.Request{}, data.User{}, raw)

	h := newVerifyTestHelper(t, pub, err)
	h.checkVerifyError(errors.New("password not matched"))
	h.checkVerifyPasswordEvent("VerifyPassword")
	h.checkVerifyPasswordFailedEvent(errors.New("password not matched"))
	h.checkAuthenticatedByPasswordEvent("")
}

// パスワードが見つからない場合、必ず失敗する
func TestVerifyFailedWhenPasswordNotFound(t *testing.T) {
	raw := data.RawPassword("password")

	pub := newVerifyTestEventPublisher()
	db := emptyVerifyTestDB()
	matcher := newVerifyTestMatcher()

	// verify
	verifier := NewVerifier(pub, db, matcher)
	err := verifier.verify(data.Request{}, data.User{}, raw)

	h := newVerifyTestHelper(t, pub, err)
	h.checkVerifyError(errors.New("password not found"))
	h.checkVerifyPasswordEvent("VerifyPassword")
	h.checkVerifyPasswordFailedEvent(errors.New("password not found"))
	h.checkAuthenticatedByPasswordEvent("")
}

type (
	verifyTestEventPublisher struct {
		verifyPassword          string
		verifyPasswordFailed    error
		authenticatedByPassword string
	}

	verifyTestDB struct {
		password *data.HashedPassword
	}

	verifyTestMatcher struct {
	}

	verifyTestHelper struct {
		t   *testing.T
		pub *verifyTestEventPublisher
		err error
	}
)

func newVerifyTestEventPublisher() *verifyTestEventPublisher {
	return &verifyTestEventPublisher{}
}

func (pub *verifyTestEventPublisher) VerifyPassword(request data.Request, user data.User) {
	pub.verifyPassword = "VerifyPassword"
}
func (pub *verifyTestEventPublisher) VerifyPasswordFailed(request data.Request, user data.User, err error) {
	pub.verifyPasswordFailed = err
}
func (pub *verifyTestEventPublisher) AuthenticatedByPassword(request data.Request, user data.User) {
	pub.authenticatedByPassword = "AuthenticatedByPassword"
}

func newVerifyTestDB(password data.HashedPassword) verifyTestDB {
	return verifyTestDB{
		password: &password,
	}
}

func emptyVerifyTestDB() verifyTestDB {
	return verifyTestDB{}
}

func (db verifyTestDB) FindUserPassword(user data.User) (data.HashedPassword, error) {
	if db.password == nil {
		return nil, errors.New("password not found")
	}
	return *db.password, nil
}

func newVerifyTestMatcher() verifyTestMatcher {
	return verifyTestMatcher{}
}

func (verifyTestMatcher) MatchPassword(hashed data.HashedPassword, raw data.RawPassword) error {
	if string(raw) != string(hashed) {
		return errors.New("password not matched")
	}
	return nil
}

func newVerifyTestHelper(t *testing.T, pub *verifyTestEventPublisher, err error) verifyTestHelper {
	return verifyTestHelper{
		t:   t,
		pub: pub,
		err: err,
	}
}

func (h verifyTestHelper) checkVerifyError(err error) {
	if err == nil {
		if h.err != nil {
			h.t.Errorf("verify fired: %s", h.err)
		}
	} else {
		if h.err.Error() != err.Error() {
			h.t.Errorf("verify error message is not matched: %s (expected: %s)", h.err, err)
		}
	}
}
func (h verifyTestHelper) checkVerifyPasswordEvent(event string) {
	if h.pub.verifyPassword != event {
		h.t.Errorf("VerifyPassword event not match: %s (expected: %s)", h.pub.verifyPassword, event)
	}
}
func (h verifyTestHelper) checkVerifyPasswordFailedEvent(err error) {
	if err == nil {
		if h.pub.verifyPasswordFailed != nil {
			h.t.Errorf("VerifyPasswordFailed event fired: %s", h.pub.verifyPasswordFailed)
		}
	} else {
		if h.pub.verifyPasswordFailed == nil {
			h.t.Error("VerifyPasswordFailed event not fired")
		}
		if h.pub.verifyPasswordFailed.Error() != err.Error() {
			h.t.Errorf("VerifyPasswordFailed error message is not matched: %s (expected: %s)", h.pub.verifyPasswordFailed, err)
		}
	}
}
func (h verifyTestHelper) checkAuthenticatedByPasswordEvent(event string) {
	if h.pub.authenticatedByPassword != event {
		h.t.Errorf("AuthenticatedByPassword event not match: %s (expected: %s)", h.pub.authenticatedByPassword, event)
	}
}
