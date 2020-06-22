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
	Time              string                `json:"time"`
	Level             string                `json:"level"`
	Message           string                `json:"message"`
	Request           *RequestLog           `json:"request"`
	AuthenticatedUser *AuthenticatedUserLog `json:"authenticated_user,omitempty"`
	Profile           *ProfileLog           `json:"profile,omitempty"`
	Resource          *ResourceLog          `json:"resource,omitempty"`
	Error             string                `json:"error,omitempty"`
}

type RequestLog struct {
	RequestedAt string   `json:"requested_at"`
	Route       RouteLog `json:"route"`
}
type AuthenticatedUserLog struct {
	UserID string `json:"user_id"`
}
type ProfileLog struct {
	Roles []string `json:"roles"`
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

func (logger Logger) Debugf(request data.Request, format string, v ...interface{}) {
	logger.logger.Debug(jsonMessage("DEBUG", Log{
		Message: fmt.Sprintf(format, v...),
		Request: requestLog(&request),
	}))
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

	if log.Request != nil {
		entry.Request = requestLog(log.Request)
	}

	if log.AuthenticatedUser != nil {
		entry.AuthenticatedUser = userLog(log.AuthenticatedUser)
	}

	if log.Profile != nil {
		entry.Profile = profileLog(log.Profile)
	}

	if log.Resource != nil {
		entry.Resource = resourceLog(log.Resource)
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

func requestLog(request *data.Request) *RequestLog {
	return &RequestLog{
		RequestedAt: request.RequestedAt.String(),
		Route: RouteLog{
			RemoteAddr: string(request.Route.RemoteAddr),
		},
	}
}

func userLog(user *data.AuthenticatedUser) *AuthenticatedUserLog {
	return &AuthenticatedUserLog{
		UserID: string(user.UserID),
	}
}

func profileLog(profile *data.Profile) *ProfileLog {
	return &ProfileLog{
		Roles: profile.Roles,
	}
}

func resourceLog(resource *data.Resource) *ResourceLog {
	return &ResourceLog{
		Path: string(resource.Path),
	}
}
