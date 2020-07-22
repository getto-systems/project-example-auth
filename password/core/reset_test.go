package core

import (
	"fmt"
	//"strings"
	"time"

	"github.com/getto-systems/project-example-id/password/db"
	"github.com/getto-systems/project-example-id/password/log"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

// 再設定トークン発行
func Example_issueResetToken() {
	h := newResetTestHelper()
	logger, db, exp, gen, testLogger := h.setup()
	h.registerLogin(db) // ログインID を登録

	request, login, _ := h.context()

	resetter := newResetter(logger, db, exp, gen)
	reset, err := resetter.issueReset(request, login)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("reset: %s\n", formatReset(&reset))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: nil
	// reset: {reset-id}
	// debug: ["Password/Reset/TryToIssueReset", req: {reset-remote}, login: {reset-login}, reset: nil, user: nil, expires: "2020-01-01 01:00:00 +0000 UTC", err: nil]
	// info: []
	// audit: ["Password/Reset/IssuedReset", req: {reset-remote}, login: {reset-login}, reset: {reset-id}, user: {reset-user}, expires: "2020-01-01 01:00:00 +0000 UTC", err: nil]
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

func (resetTestGenerator) Generate() (password.ResetID, password.ResetToken, error) {
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

func (h resetTestHelper) setup() (password.Logger, *db.MemoryStore, password.Expiration, password.ResetGenerator, *testLogger) {
	testLogger := newTestLogger()
	logger := log.NewLogger(testLogger)

	db := db.NewMemoryStore()

	exp := password.NewExpiration(data.Hour(1))

	return logger, db, exp, h.gen, testLogger
}

func (h resetTestHelper) registerLogin(db *db.MemoryStore) {
	db.RegisterLogin(h.user, h.login)
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
		formatReset(entry.Reset),
		formatUser(entry.User),
		formatExpires(entry.Expires),
		formatError(entry.Error),
	)
}

func (h resetTestHelper) formatDB(db *db.MemoryStore) string {
	resetID, _, _ := h.gen.Generate()
	reset, ok := db.GetResetUser(resetID)
	if !ok {
		return "nil"
	} else {
		return fmt.Sprintf("\"%s\"", reset.User().UserID())
	}
}
