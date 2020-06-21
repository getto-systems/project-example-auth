package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/getto-systems/project-example-id/user/subscriber"

	"github.com/getto-systems/project-example-id/data"
)

type Log struct {
	Time     string      `json:"time"`
	Level    string      `json:"level"`
	Message  string      `json:"message"`
	Request  RequestLog  `json:"request"`
	UserID   string      `json:"user_id,omitempty"`
	Resource ResourceLog `json:"resource,omitempty"`
	Error    string      `json:"error,omitempty"`
}

type ResourceLog struct {
	Path string `json:"path"`
}
type RouteLog struct {
	RemoteAddr string `json:"remote_addr"`
}
type RequestLog struct {
	RequestedAt string   `json:"requested_at"`
	Route       RouteLog `json:"route"`
}

type Logger struct {
	logger *log.Logger
}

func (logger Logger) Audit(log subscriber.Log) {
	logger.json("AUDIT", format(log))
}

func (logger Logger) Info(log subscriber.Log) {
	logger.json("INFO", format(log))
}

func (logger Logger) Debug(log subscriber.Log) {
	logger.json("DEBUG", format(log))
}

func (logger Logger) Debugf(request data.Request, format string, v ...interface{}) {
	logger.json("DEBUG", Log{
		Message: fmt.Sprintf(format, v...),
		Request: requestLog(request),
	})
}

func (logger Logger) json(level string, log Log) {
	logger.logger.Println(jsonMessage(level, log))
}

func NewLogger(level string, logger *log.Logger) Logger {
	return Logger{
		logger: leveledLogger(level, logger),
	}
}
func leveledLogger(level string, logger *log.Logger) *log.Logger {
	// TODO applog に差し替え
	return logger
}

func format(log subscriber.Log) Log {
	return Log{
		Time:     time.Now().UTC().String(),
		Message:  log.Message,
		Request:  requestLog(log.Request),
		UserID:   string(log.UserID),
		Resource: resourceLog(log.Resource),
		Error:    log.Error.Error(),
	}
}
func jsonMessage(level string, log Log) []byte {
	log.Level = level
	data, err := json.Marshal(log)
	if err != nil {
		return nil
	}

	return data
}

func requestLog(request data.Request) RequestLog {
	return RequestLog{
		RequestedAt: request.RequestedAt.String(),
		Route: RouteLog{
			RemoteAddr: string(request.Route.RemoteAddr),
		},
	}
}

func resourceLog(resource data.Resource) ResourceLog {
	return ResourceLog{
		Path: string(resource.Path),
	}
}
