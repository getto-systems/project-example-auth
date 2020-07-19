package core

import (
	"github.com/getto-systems/project-example-id/event_log"

	"testing"
)

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

func assertError(t *testing.T, label string, got error, expected error) {
	if expected == nil {
		if got != nil {
			t.Errorf("%s error: %s", label, got)
		}
	} else {
		if got == nil {
			t.Errorf("%s success", label)
		} else {
			if got.Error() != expected.Error() {
				t.Errorf("%s error message is not matched: %s (expected: %s)", label, got, expected)
			}
		}
	}
}
