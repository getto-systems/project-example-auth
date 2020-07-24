package core

import (
	"fmt"
	"strings"

	"github.com/getto-systems/project-example-id/password/log"
	repository_password "github.com/getto-systems/project-example-id/password/repository/password"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

// パスワードが一致したら audit: authenticated by password
func Example_validate() {
	h := newValidateTestHelper()
	logger, passwords, matcher, testLogger := h.setup()
	h.registerPassword(passwords, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword("password") // 保存されているものと同じパスワード

	validator := newValidator(logger, passwords, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: nil
	// user: {validate-user}
	// debug: ["Password/Validate/TryToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["Password/Validate/AuthedByPassword", req: {validate-remote}, login: {validate-login}, user: {validate-user}, err: nil]
}

// パスワードが一致しなかったら audit: validate password failed
func Example_validate_fail_DifferentPassword() {
	h := newValidateTestHelper()
	logger, passwords, matcher, testLogger := h.setup()
	h.registerPassword(passwords, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword("different-password") // 保存されているものと違うパスワード

	validator := newValidator(logger, passwords, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: "Password/Password/NotMatched"
	// user: {}
	// debug: ["Password/Validate/TryToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["Password/Validate/FailedToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: "Password/Password/NotMatched"]
}

// パスワードが見つからない場合は認証失敗
func Example_validate_fail_PasswordNotFound() {
	h := newValidateTestHelper()
	logger, passwords, matcher, testLogger := h.setup()
	//h.registerPassword(passwords, password.HashedPassword("password")) // パスワードの登録はしない

	request, login := h.context()
	raw := password.RawPassword("password") // このユーザーのパスワードは登録されていない

	validator := newValidator(logger, passwords, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: "Password/Password/NotFound/Password"
	// user: {}
	// debug: ["Password/Validate/TryToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["Password/Validate/FailedToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: "Password/Password/NotFound/Password"]
}

// 空のパスワードの場合、必ず失敗する
func Example_validate_fail_EmptyPassword() {
	h := newValidateTestHelper()
	logger, passwords, matcher, testLogger := h.setup()
	h.registerPassword(passwords, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword("") // 空のパスワード

	validator := newValidator(logger, passwords, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: "Password/Password/Empty"
	// user: {}
	// debug: ["Password/Validate/TryToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["Password/Validate/FailedToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: "Password/Password/Empty"]
}

// 長いパスワードの場合、必ず失敗する
func Example_validate_fail_LongPassword() {
	h := newValidateTestHelper()
	logger, passwords, matcher, testLogger := h.setup()
	h.registerPassword(passwords, password.HashedPassword("password"))

	request, login := h.context()
	raw := password.RawPassword(strings.Repeat("a", 73)) // 長いパスワード

	validator := newValidator(logger, passwords, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: "Password/Password/TooLong"
	// user: {}
	// debug: ["Password/Validate/TryToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["Password/Validate/FailedToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: "Password/Password/TooLong"]
}

// ギリギリの長さのパスワードの場合、成功する
func Example_validate_LongPassword() {
	h := newValidateTestHelper()
	logger, passwords, matcher, testLogger := h.setup()
	h.registerPassword(passwords, password.HashedPassword(strings.Repeat("a", 72)))

	request, login := h.context()
	raw := password.RawPassword(strings.Repeat("a", 72)) // 72 バイトまで許容

	validator := newValidator(logger, passwords, matcher)
	user, err := validator.validate(request, login, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("user: %s\n", formatUser(&user))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: nil
	// user: {validate-user}
	// debug: ["Password/Validate/TryToValidate", req: {validate-remote}, login: {validate-login}, user: nil, err: nil]
	// info: []
	// audit: ["Password/Validate/AuthedByPassword", req: {validate-remote}, login: {validate-login}, user: {validate-user}, err: nil]
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

func newValidateTestMatcher() (matcher validateTestMatcher) {
	return
}

func (validateTestMatcher) MatchPassword(hashed password.HashedPassword, raw password.RawPassword) error {
	if string(raw) != string(hashed) {
		return password.ErrPasswordNotMatched
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

func (h validateTestHelper) setup() (password.Logger, *repository_password.MemoryStore, password.PasswordMatcher, *testLogger) {
	testLogger := newTestLogger()
	logger := log.NewLogger(testLogger)

	passwords := repository_password.NewMemoryStore()

	return logger, passwords, h.matcher, testLogger
}

func (h validateTestHelper) registerPassword(passwords *repository_password.MemoryStore, password password.HashedPassword) {
	passwords.RegisterPassword(h.user, password)
	passwords.RegisterLogin(h.user, h.login)
}

func (h validateTestHelper) context() (data.Request, password.Login) {
	return h.request, h.login
}

func (h validateTestHelper) formatLog(entry event_log.Entry) string {
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
