package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/getto-systems/applog-go/v2"

	"github.com/getto-systems/project-example-id/user/subscriber"

	"github.com/getto-systems/project-example-id/data"
)

type Log struct {
	Time    string     `json:"time"`
	Level   string     `json:"level"`
	Message string     `json:"message"`
	Request RequestLog `json:"request"`
	User    *UserLog   `json:"user,omitempty"`
	Ticket  *TicketLog `json:"ticket,omitempty"`
	Error   string     `json:"error,omitempty"`
}

type RequestLog struct {
	RequestedAt string   `json:"requested_at"`
	Route       RouteLog `json:"route"`
}
type UserLog struct {
	UserID string `json:"user_id"`
}
type TicketLog struct {
	Profile         ProfileLog `json:"profile"`
	AuthenticatedAt string     `json:"authenticated_at"`
	Expires         string     `json:"expires"`
}
type ProfileLog struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
}
type ResourceLog struct {
	Path string `json:"path"`
}
type RouteLog struct {
	RemoteAddr string `json:"remote_addr"`
}

type Logger struct {
	logger applog.Logger
}

func (logger Logger) Audit(entry subscriber.Log) {
	logger.logger.Always(jsonMessage("AUDIT", format(entry)))
}

func (logger Logger) Info(entry subscriber.Log) {
	logger.logger.Info(jsonMessage("INFO", format(entry)))
}

func (logger Logger) Debug(entry subscriber.Log) {
	logger.logger.Debug(jsonMessage("DEBUG", format(entry)))
}

func (logger Logger) DebugMessage(request data.Request, message string) {
	logger.logger.Debug(jsonMessage("DEBUG", Log{
		Message: message,
		Request: requestLog(request),
	}))
}

func (logger Logger) DebugError(request data.Request, format string, err error) {
	logger.DebugMessage(request, fmt.Sprintf(format, err))
}

func NewLogger(level string, logger *log.Logger) Logger {
	return Logger{
		logger: leveledLogger(level, logger),
	}
}
func leveledLogger(level string, logger *log.Logger) applog.Logger {
	switch level {
	case "QUIET":
		return applog.NewQuietLogger(logger)
	case "INFO":
		return applog.NewInfoLogger(logger)
	default:
		return applog.NewDebugLogger(logger)
	}
}

func format(log subscriber.Log) Log {
	entry := Log{
		Time:    time.Now().UTC().String(),
		Message: log.Message,
	}

	entry.Request = requestLog(log.Request)

	if log.User != nil {
		entry.User = userLog(log.User)
	}

	if log.Ticket != nil {
		entry.Ticket = ticketLog(log.Ticket)
	}

	if log.Error != nil {
		entry.Error = log.Error.Error()
	}

	return entry
}
func jsonMessage(level string, log Log) string {
	log.Level = level
	data, err := json.Marshal(log)
	if err != nil {
		return "json marshal error"
	}

	return string(data)
}

func requestLog(request data.Request) RequestLog {
	return RequestLog{
		RequestedAt: request.RequestedAt.String(),
		Route: RouteLog{
			RemoteAddr: string(request.Route.RemoteAddr),
		},
	}
}

func userLog(user *data.User) *UserLog {
	return &UserLog{
		UserID: string(user.UserID),
	}
}

func ticketLog(ticket *data.Ticket) *TicketLog {
	return &TicketLog{
		Profile:         profileLog(ticket.Profile),
		AuthenticatedAt: ticket.AuthenticatedAt.String(),
		Expires:         ticket.Expires.String(),
	}
}

func profileLog(profile data.Profile) ProfileLog {
	return ProfileLog{
		UserID: string(profile.UserID),
		Roles:  profile.Roles,
	}
}
