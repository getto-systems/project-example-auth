package core

import (
	"fmt"
	//"strings"
	"time"

	"github.com/getto-systems/project-example-id/password/log"
	repository_password "github.com/getto-systems/project-example-id/password/repository/password"
	repository_session "github.com/getto-systems/project-example-id/password/repository/reset_session"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

// 再設定トークン発行
func Example_issueResetToken() {
	h := newResetTestHelper()
	logger, passwords, sessions, exp, gen, testLogger := h.setup()
	h.registerLogin(passwords) // ログインID を登録

	request, login, _ := h.context()

	resetter := newResetter(logger, passwords, sessions, exp, gen)
	reset, err := resetter.createResetSession(request, login)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("reset: %s\n", formatResetSession(&reset))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: nil
	// reset: {reset-id}
	// debug: ["Password/Reset/TryToCreateResetSession", req: {reset-remote}, login: {reset-login}, reset: nil, user: nil, expires: "2020-01-01 01:00:00 +0000 UTC", err: nil]
	// info: []
	// audit: ["Password/Reset/CreatedResetSession", req: {reset-remote}, login: {reset-login}, reset: {reset-id}, user: {reset-user}, expires: "2020-01-01 01:00:00 +0000 UTC", err: nil]
}

// ログインID が登録されていない場合は発行できない

// 再設定ステータスを取得
// 登録されていない場合は「リクエストなし」

// 再設定トークン検証
// 再設定トークン検証

type (
	resetTestGenerator struct{}

	resetTestHelper struct {
		gen resetTestGenerator

		request data.Request
		user    data.User
		login   password.Login
	}

	resetTestLogEntry struct {
		message string
		err     error
	}
)

func newResetTestGenerator() resetTestGenerator {
	return resetTestGenerator{}
}

func (resetTestGenerator) GenerateSession() (password.ResetSessionID, password.ResetToken, error) {
	return "reset-id", "reset-token", nil
}

func newResetTestHelper() resetTestHelper {
	gen := newResetTestGenerator()

	now, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	request := data.NewRequest(data.RequestedAt(now), data.RemoteAddr("reset-remote"))
	user := data.NewUser("reset-user")
	login := password.NewLogin("reset-login")

	return resetTestHelper{
		gen: gen,

		request: request,
		user:    user,
		login:   login,
	}
}

func (h resetTestHelper) setup() (password.Logger, *repository_password.MemoryStore, *repository_session.MemoryStore, password.ResetSessionExpiration, password.ResetSessionGenerator, *testLogger) {
	testLogger := newTestLogger()
	logger := log.NewLogger(testLogger)

	passwords := repository_password.NewMemoryStore()
	sessions := repository_session.NewMemoryStore()

	exp := password.NewResetSessionExpiration(data.Hour(1))

	return logger, passwords, sessions, exp, h.gen, testLogger
}

func (h resetTestHelper) registerLogin(passwords *repository_password.MemoryStore) {
	passwords.RegisterLogin(h.user, h.login)
}

func (h resetTestHelper) context() (data.Request, password.Login, data.User) {
	return h.request, h.login, h.user
}

func (h resetTestHelper) formatLog(entry event_log.Entry) string {
	if entry.Message == "" {
		return "[]"
	}

	return fmt.Sprintf(
		"[\"%s\", req: %s, login: %s, reset: %s, user: %s, expires: %s, err: %s]",
		entry.Message,
		formatRequest(entry.Request),
		formatLogin(entry.Login),
		formatResetSession(entry.ResetSession),
		formatUser(entry.User),
		formatExpires(entry.Expires),
		formatError(entry.Error),
	)
}

func (h resetTestHelper) formatDB(sessions *repository_session.MemoryStore) string {
	id, _, _ := h.gen.GenerateSession()
	session, ok := sessions.GetResetSessionData(id)
	if !ok {
		return "nil"
	} else {
		return fmt.Sprintf("\"%s\"", session.User().UserID())
	}
}
