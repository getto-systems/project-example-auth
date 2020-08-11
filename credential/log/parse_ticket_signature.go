package credential_log

import (
	"fmt"

	"github.com/getto-systems/project-example-id/credential/infra"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func (log Logger) parseTicketSignature() infra.ParseTicketSignatureLogger {
	return log
}

func (log Logger) TryToParseTicketSignature(request request.Request) {
	log.logger.Debug(parseTicketSignatureLog("TryToParse", request, nil, nil))
}
func (log Logger) FailedToParseTicketSignature(request request.Request, err error) {
	log.logger.Error(parseTicketSignatureLog("FailedToParse", request, nil, err))
}
func (log Logger) FailedToParseTicketSignatureBecauseNonceMatchFailed(request request.Request, err error) {
	log.logger.Audit(parseTicketSignatureLog("FailedToParseBecauseNonceMatchFailed", request, nil, err))
}
func (log Logger) ParseTicketSignature(request request.Request, user user.User) {
	log.logger.Info(parseTicketSignatureLog("Parse", request, &user, nil))
}

type (
	parseTicketSignatureEntry struct {
		Action  string             `json:"action"`
		Message string             `json:"message"`
		Request request.RequestLog `json:"request"`
		User    *user.UserLog      `json:"user,omitempty"`
		Err     *string            `json:"error,omitempty"`
	}
)

func parseTicketSignatureLog(message string, request request.Request, user *user.User, err error) parseTicketSignatureEntry {
	entry := parseTicketSignatureEntry{
		Action:  "Credential/ParseTicketSignature",
		Message: message,
		Request: request.Log(),
	}

	if user != nil {
		log := user.Log()
		entry.User = &log
	}
	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry parseTicketSignatureEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
