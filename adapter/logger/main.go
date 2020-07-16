package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/getto-systems/applog-go/v2"

	"github.com/getto-systems/project-example-id/http_handler"

	"github.com/getto-systems/project-example-id/event_log"

	"github.com/getto-systems/project-example-id/ticket"

	"github.com/getto-systems/project-example-id/data"
)

type Logger struct {
	logger applog.Logger
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

func (logger Logger) Logger() event_log.Logger {
	return logger
}

func (logger Logger) RequestLogger() http_handler.RequestLogger {
	return logger
}

func (logger Logger) Audit(entry event_log.Entry) {
	logger.logger.Always(jsonMessage("AUDIT", format(entry)))
}

func (logger Logger) Info(entry event_log.Entry) {
	logger.logger.Info(jsonMessage("INFO", format(entry)))
}

func (logger Logger) Debug(entry event_log.Entry) {
	logger.logger.Debug(jsonMessage("DEBUG", format(entry)))
}

func (logger Logger) DebugMessage(request data.Request, message string) {
	logger.logger.Debug(jsonMessage("DEBUG", Entry{
		Message: message,
		Request: requestLog(request),
	}))
}

func (logger Logger) DebugError(request data.Request, format string, err error) {
	logger.DebugMessage(request, fmt.Sprintf(format, err))
}

type Entry struct {
	Time    string      `json:"time"`
	Level   string      `json:"level"`
	Message string      `json:"message"`
	Request RequestLog  `json:"request"`
	Nonce   *NonceLog   `json:"nonce,omitempty"`
	User    *UserLog    `json:"user,omitempty"`
	Roles   *RolesLog   `json:"roles,omitempty"`
	Expires *ExpiresLog `json:"expires,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type RequestLog struct {
	RequestedAt string   `json:"requested_at"`
	Route       RouteLog `json:"route"`
}
type RouteLog struct {
	RemoteAddr string `json:"remote_addr"`
}
type NonceLog struct {
	Nonce string `json:"nonce"`
}
type UserLog struct {
	UserID string `json:"user_id"`
}
type RolesLog struct {
	Roles []string `json:"roles"`
}
type ExpiresLog struct {
	Expires string `json:"expires"`
}

func format(log event_log.Entry) Entry {
	entry := Entry{
		Time:    time.Now().UTC().String(),
		Message: log.Message,
	}

	entry.Request = requestLog(log.Request)

	if log.Nonce != nil {
		entry.Nonce = nonceLog(log.Nonce)
	}

	if log.User != nil {
		entry.User = userLog(log.User)
	}

	if log.Roles != nil {
		entry.Roles = rolesLog(log.Roles)
	}

	if log.Expires != nil {
		entry.Expires = expiresLog(log.Expires)
	}

	if log.Error != nil {
		entry.Error = log.Error.Error()
	}

	return entry
}
func jsonMessage(level string, log Entry) string {
	log.Level = level
	data, err := json.Marshal(log)
	if err != nil {
		return "json marshal error"
	}

	return string(data)
}

func requestLog(request data.Request) RequestLog {
	return RequestLog{
		RequestedAt: time.Time(request.RequestedAt()).String(),
		Route: RouteLog{
			RemoteAddr: string(request.Route().RemoteAddr()),
		},
	}
}

func nonceLog(nonce *ticket.Nonce) *NonceLog {
	return &NonceLog{
		Nonce: string(*nonce),
	}
}

func userLog(user *data.User) *UserLog {
	return &UserLog{
		UserID: string(user.UserID()),
	}
}

func rolesLog(roles *data.Roles) *RolesLog {
	log := RolesLog{}
	for _, role := range *roles {
		log.Roles = append(log.Roles, string(role))
	}

	return &log
}

func expiresLog(expires *data.Expires) *ExpiresLog {
	return &ExpiresLog{
		Expires: time.Time(*expires).String(),
	}
}
