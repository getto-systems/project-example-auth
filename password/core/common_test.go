package core

import (
	"fmt"
	"time"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

type testReportError func(string)

type testLogger struct {
	audit event_log.Entry
	info  event_log.Entry
	debug event_log.Entry
}

func newTestLogger() *testLogger {
	return &testLogger{}
}

func (logger *testLogger) eventLogger() event_log.Logger {
	return logger
}

func (logger *testLogger) Audit(entry event_log.Entry) {
	logger.audit = entry
}
func (logger *testLogger) Info(entry event_log.Entry) {
	logger.info = entry
}
func (logger *testLogger) Debug(entry event_log.Entry) {
	logger.debug = entry
}

func formatError(err error, expected error) string {
	if err == nil {
		return "err: nil"
	}

	if err != expected {
		return fmt.Sprintf("err: \"%s\"", err)
	}

	return "err"
}

func formatRequest(request data.Request, expected data.Request) string {
	if request.Route().RemoteAddr() != expected.Route().RemoteAddr() {
		return fmt.Sprintf("req: {%s}", request.Route().RemoteAddr())
	}

	return "req"
}

func formatLogin(login *password.Login, expected *password.Login) string {
	if login == nil {
		return "login: nil"
	}

	if login.ID() != expected.ID() {
		return fmt.Sprintf("login: {%s}", login.ID())
	}

	return "login"
}

func formatResetSession(session *password.ResetSession) string {
	if session == nil {
		return "session: nil"
	}

	return fmt.Sprintf("session: {%s}", session.ID())
}

func formatUser(user *data.User, expected *data.User) string {
	if user == nil {
		return "user: nil"
	}

	if user.UserID() != expected.UserID() {
		return fmt.Sprintf("user: {%s}", user.UserID())
	}

	return "user"
}

func formatExpires(expires *data.Expires, expected *data.Expires) string {
	if expires == nil {
		return "expires: nil"
	}

	if expected == nil || *expires != *expected {
		return fmt.Sprintf("expires: \"%s\"", time.Time(*expires).String())
	}

	return "expires"
}
