package ticket_log

import (
	"fmt"
	"time"

	"github.com/getto-systems/project-example-auth/ticket/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

func (log Logger) register() infra.RegisterLogger {
	return log
}

func (log Logger) TryToRegister(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit) {
	log.logger.Debug(registerLog("TryToRegister", request, user, expires, limit, nil))
}
func (log Logger) FailedToRegister(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit, err error) {
	log.logger.Error(registerLog("FailedToRegister", request, user, expires, limit, err))
}
func (log Logger) Register(request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit) {
	log.logger.Info(registerLog("Register", request, user, expires, limit, nil))
}

type (
	registerEntry struct {
		Action      string             `json:"action"`
		Message     string             `json:"message"`
		Request     request.RequestLog `json:"request"`
		User        user.UserLog       `json:"user"`
		Expires     string             `json:"expires"`
		ExtendLimit string             `json:"extend_limit"`
		Err         *string            `json:"error,omitempty"`
	}
)

func registerLog(message string, request request.Request, user user.User, expires credential.TicketExpires, limit credential.TicketExtendLimit, err error) registerEntry {
	entry := registerEntry{
		Action:      "Ticket/Register",
		Message:     message,
		Request:     request.Log(),
		User:        user.Log(),
		Expires:     time.Time(expires).String(),
		ExtendLimit: time.Time(limit).String(),
	}

	if err != nil {
		message := err.Error()
		entry.Err = &message
	}

	return entry
}

func (entry registerEntry) String() string {
	return fmt.Sprintf("%s/%s", entry.Action, entry.Message)
}
