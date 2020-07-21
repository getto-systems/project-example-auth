package core

import (
	"fmt"

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

func formatError(err error) string {
	if err == nil {
		return "nil"
	} else {
		return fmt.Sprintf("\"%s\"", err)
	}
}

func formatRequest(request data.Request) string {
	return fmt.Sprintf("{%s}", request.Route().RemoteAddr())
}

func formatLogin(login *password.Login) string {
	if login == nil {
		return "nil"
	} else {
		return fmt.Sprintf("{%s}", login.ID())
	}
}

func formatUser(user *data.User) string {
	if user == nil {
		return "nil"
	} else {
		return fmt.Sprintf("{%s}", user.UserID())
	}
}
