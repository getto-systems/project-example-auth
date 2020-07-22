package core

import (
	"fmt"
	"strings"

	"github.com/getto-systems/project-example-id/password/db"
	password_event_log "github.com/getto-systems/project-example-id/password/event_log"
	"github.com/getto-systems/project-example-id/password/pubsub"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"

	"errors"
)

// パスワードが一致したら audit: authenticated by password
func Example_validate() {
	h := newValidateTestHelper()
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword("password") // 保存されているものと同じパスワード

	validator := newValidator(pub, db, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", formatValidateLog(logger.debug))
	fmt.Printf("info: %s\n", formatValidateLog(logger.info))
	fmt.Printf("audit: %s\n", formatValidateLog(logger.audit))

	// Output:
	// err: nil
	// user: {validate-user}
	// debug: ["validate password", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["authenticated by password", req: {validate-remote}, login: {validate-login}, user: {validate-user}, err: nil]
}

// パスワードが一致しなかったら audit: validate password failed
func Example_validate_fail_DifferentPassword() {
	h := newValidateTestHelper()
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword("different-password") // 保存されているものと違うパスワード

	validator := newValidator(pub, db, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", formatValidateLog(logger.debug))
	fmt.Printf("info: %s\n", formatValidateLog(logger.info))
	fmt.Printf("audit: %s\n", formatValidateLog(logger.audit))

	// Output:
	// err: "password not matched"
	// user: {}
	// debug: ["validate password", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["validate password failed", req: {validate-remote}, login: {validate-login}, user: nil, err: "password not matched"]
}

// パスワードが見つからない場合は認証失敗
func Example_validate_fail_PasswordNotFound() {
	h := newValidateTestHelper()
	pub, db, matcher, logger := h.setup()
	//h.registerPassword(db, password.HashedPassword("password")) // パスワードの登録はしない

	request, login := h.context()
	raw := password.RawPassword("password") // このユーザーのパスワードは登録されていない

	validator := newValidator(pub, db, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", formatValidateLog(logger.debug))
	fmt.Printf("info: %s\n", formatValidateLog(logger.info))
	fmt.Printf("audit: %s\n", formatValidateLog(logger.audit))

	// Output:
	// err: "password not found"
	// user: {}
	// debug: ["validate password", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["validate password failed", req: {validate-remote}, login: {validate-login}, user: nil, err: "password not found"]
}

// 空のパスワードの場合、必ず失敗する
func Example_validate_fail_EmptyPassword() {
	h := newValidateTestHelper()
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword("") // 空のパスワード

	validator := newValidator(pub, db, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", formatValidateLog(logger.debug))
	fmt.Printf("info: %s\n", formatValidateLog(logger.info))
	fmt.Printf("audit: %s\n", formatValidateLog(logger.audit))

	// Output:
	// err: "password is empty"
	// user: {}
	// debug: ["validate password", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["validate password failed", req: {validate-remote}, login: {validate-login}, user: nil, err: "password is empty"]
}

// 長いパスワードの場合、必ず失敗する
func Example_validate_fail_LongPassword() {
	h := newValidateTestHelper()
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword(strings.Repeat("a", 73)) // 長いパスワード

	validator := newValidator(pub, db, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", formatValidateLog(logger.debug))
	fmt.Printf("info: %s\n", formatValidateLog(logger.info))
	fmt.Printf("audit: %s\n", formatValidateLog(logger.audit))

	// Output:
	// err: "password is too long"
	// user: {}
	// debug: ["validate password", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["validate password failed", req: {validate-remote}, login: {validate-login}, user: nil, err: "password is too long"]
}

// ギリギリの長さのパスワードの場合、成功する
func Example_validate_LongPassword() {
	h := newValidateTestHelper()
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword(strings.Repeat("a", 72)))

	request, login := h.context()
	raw := password.RawPassword(strings.Repeat("a", 72)) // 72 バイトまで許容

	validator := newValidator(pub, db, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", formatValidateLog(logger.debug))
	fmt.Printf("info: %s\n", formatValidateLog(logger.info))
	fmt.Printf("audit: %s\n", formatValidateLog(logger.audit))

	// Output:
	// err: nil
	// user: {validate-user}
	// debug: ["validate password", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["authenticated by password", req: {validate-remote}, login: {validate-login}, user: {validate-user}, err: nil]
}

type (
	validateTestMatcher struct{}

	validateTestHelper struct {
		matcher validateTestMatcher

		request data.Request
		user    data.User
		login   password.Login
	}

	validateTestLogEntry struct {
		message string
		err     error
	}
)

func newValidateTestMatcher() validateTestMatcher {
	return validateTestMatcher{}
}

func (validateTestMatcher) MatchPassword(hashed password.HashedPassword, raw password.RawPassword) error {
	if string(raw) != string(hashed) {
		return errors.New("password not matched")
	}
	return nil
}

func newValidateTestHelper() validateTestHelper {
	matcher := newValidateTestMatcher()

	request := data.NewRequest(data.RequestedAt{}, data.RemoteAddr("validate-remote"))
	user := data.NewUser("validate-user")
	login := password.NewLogin("validate-login")

	return validateTestHelper{
		matcher: matcher,

		request: request,
		user:    user,
		login:   login,
	}
}

func (h validateTestHelper) setup() (password.EventPublisher, *db.MemoryStore, password.Matcher, *testLogger) {
	pub := pubsub.NewPubSub()
	logger := newTestLogger()
	log := password_event_log.NewEventLogger(logger)
	pub.Subscribe(log)

	db := db.NewMemoryStore()

	return pub, db, h.matcher, logger
}

func (h validateTestHelper) registerPassword(db *db.MemoryStore, password password.HashedPassword) {
	db.RegisterPassword(h.user, password)
	db.RegisterLogin(h.user, h.login)
}

func (h validateTestHelper) context() (data.Request, password.Login) {
	return h.request, h.login
}

func formatValidateLog(entry event_log.Entry) string {
	if entry.Message == "" {
		return "[]"
	}

	return fmt.Sprintf(
		"[\"%s\", req: %s, login: %s, user: %s, err: %s]",
		entry.Message,
		formatRequest(entry.Request),
		formatLogin(entry.Login),
		formatUser(entry.User),
		formatError(entry.Error),
	)
}
