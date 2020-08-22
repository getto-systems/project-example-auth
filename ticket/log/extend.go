package ticket_log

import (
	"fmt"
	"time"

	"github.com/getto-systems/project-example-auth/ticket/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) extend() infra.ExtendLogger {
	return log
}

func (log Logger) TryToExtend(request request.Request, user user.User) {
	log.logger.Debug(extendLog("TryToExtend", request, user, nil, nil, nil))
}
func (log Logger) FailedToExtend(request request.Request, user user.User, err error) {
	log.logger.Error(extendLog("FailedToExtend", request, user, nil, nil, err))
}
func (log Logger) Extend(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit) {
	log.logger.Info(extendLog("Extend", request, user, &expires, &limit, nil))
}

type (
	extendEntry struct {
		Action      string             `json:"action"`
		Message     string             `json:"message"`
		Request     request.RequestLog `json:"request"`
		User        user.UserLog       `json:"user"`
		Expires     *string            `json:"expires,omitempty"`
		ExtendLimit *string            `json:"extend_limit,omitempty"`
		Err         *string            `json:"error,omitempty"`
	}
)

func extendLog(message string, request request.Request, user user.User, expires *credential.TicketExpires, limit *credential.TicketExtendLimit, err error) extendEntry {
	entry := extendEntry{
		Action:  "Ticket/Extend",
		Message: message,
		Request: request.Log(),
		User:    user.Log(),
	}

	if expires != nil {
		log := time.Time(*expires).String()
		entry.Expires = &log
	}
	if limit != nil {
		log := time.Time(*limit).String()
		entry.ExtendLimit = &log
	}
	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry extendEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
