package logger

import (
	"encoding/json"
	golog "log"
	"os"
	gotime "time"

	"github.com/getto-systems/applog-go/v2"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type Logger struct {
	logger applog.Logger
}

func NewLogger(level string) Logger {
	return Logger{
		logger: leveledLogger(level),
	}
}
func leveledLogger(level string) applog.Logger {
	logger := golog.New(os.Stdout, "", 0)

	switch level {
	case "QUIET":
		return applog.NewQuietLogger(logger)
	case "INFO":
		return applog.NewInfoLogger(logger)
	default:
		return applog.NewDebugLogger(logger)
	}
}

func (logger Logger) log() log.Logger {
	return logger
}

func (logger Logger) Audit(entry log.Entry) {
	logger.logger.Always(jsonMessage("AUDIT", format(entry)))
}

func (logger Logger) Error(entry log.Entry) {
	logger.logger.Always(jsonMessage("ERROR", format(entry)))
}

func (logger Logger) Info(entry log.Entry) {
	logger.logger.Info(jsonMessage("INFO", format(entry)))
}

func (logger Logger) Debug(entry log.Entry) {
	logger.logger.Debug(jsonMessage("DEBUG", format(entry)))
}

type Entry struct {
	Time    string     `json:"time"`
	Level   string     `json:"level"`
	Message string     `json:"message"`
	Request RequestLog `json:"request"`

	User  *UserLog  `json:"user,omitempty"`
	Login *LoginLog `json:"login,omitempty"`

	TicketNonce *TicketNonceLog `json:"ticket,omitempty"`
	ApiRoles    *ApiRolesLog    `json:"api_roles,omitempty"`

	Expires     *ExpiresLog     `json:"expires,omitempty"`
	ExtendLimit *ExtendLimitLog `json:"extend_limit,omitempty"`

	ResetSession     *ResetSessionLog     `json:"reset_session,omitempty"`
	ResetStatus      *ResetStatusLog      `json:"reset_status,omitempty"`
	ResetDestination *ResetDestinationLog `json:"reset_destination,omitempty"`

	Error string `json:"error,omitempty"`
}

type RequestLog struct {
	RequestedAt string   `json:"requested_at"`
	Route       RouteLog `json:"route"`
}
type RouteLog struct {
	RemoteAddr string `json:"remote_addr"`
}

type UserLog struct {
	UserID string `json:"user_id"`
}
type LoginLog struct {
	LoginID string `json:"login_id"`
}

type TicketNonceLog struct {
	Nonce string `json:"nonce"`
}
type ApiRolesLog struct {
	ApiRoles []string `json:"api_roles"`
}

type ExpiresLog struct {
	Expires string `json:"expires"`
}
type ExtendLimitLog struct {
	ExtendLimit string `json:"extend_limit"`
}

type ResetSessionLog struct {
	SessionID string `json:"session_id"`
}

// TODO あとでちゃんとする
type ResetStatusLog struct {
}

// TODO あとでちゃんとする
type ResetDestinationLog struct {
}

func format(log log.Entry) Entry {
	entry := Entry{
		Time:    gotime.Now().UTC().String(),
		Message: log.Message,
	}

	entry.Request = requestLog(log.Request)

	if log.User != nil {
		entry.User = userLog(log.User)
	}
	if log.Login != nil {
		entry.Login = loginLog(log.Login)
	}

	if log.TicketNonce != nil {
		entry.TicketNonce = ticketNonceLog(log.TicketNonce)
	}
	if log.ApiRoles != nil {
		entry.ApiRoles = apiRolesLog(log.ApiRoles)
	}

	if log.Expires != nil {
		entry.Expires = expiresLog(log.Expires)
	}
	if log.ExtendLimit != nil {
		entry.ExtendLimit = extendLimitLog(log.ExtendLimit)
	}

	if log.ResetSession != nil {
		entry.ResetSession = resetSessionLog(log.ResetSession)
	}
	if log.ResetStatus != nil {
		entry.ResetStatus = resetStatusLog(log.ResetStatus)
	}
	if log.ResetDestination != nil {
		entry.ResetDestination = resetDestinationLog(log.ResetDestination)
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

func requestLog(request request.Request) RequestLog {
	return RequestLog{
		RequestedAt: gotime.Time(request.RequestedAt()).String(),
		Route: RouteLog{
			RemoteAddr: string(request.Route().RemoteAddr()),
		},
	}
}

func userLog(user *user.User) *UserLog {
	return &UserLog{
		UserID: string(user.ID()),
	}
}
func loginLog(login *user.Login) *LoginLog {
	return &LoginLog{
		LoginID: string(login.ID()),
	}
}

func ticketNonceLog(nonce *credential.TicketNonce) *TicketNonceLog {
	return &TicketNonceLog{
		Nonce: string(*nonce),
	}
}
func apiRolesLog(roles *credential.ApiRoles) *ApiRolesLog {
	log := ApiRolesLog{}
	for _, role := range *roles {
		log.ApiRoles = append(log.ApiRoles, string(role))
	}

	return &log
}

func expiresLog(expires *time.Expires) *ExpiresLog {
	return &ExpiresLog{
		Expires: gotime.Time(*expires).String(),
	}
}
func extendLimitLog(limit *time.ExtendLimit) *ExtendLimitLog {
	return &ExtendLimitLog{
		ExtendLimit: gotime.Time(*limit).String(),
	}
}

func resetSessionLog(session *password_reset.Session) *ResetSessionLog {
	return &ResetSessionLog{
		SessionID: string(session.ID()),
	}
}
func resetStatusLog(status *password_reset.Status) *ResetStatusLog {
	// TODO あとでちゃんとする
	return &ResetStatusLog{}
}
func resetDestinationLog(dest *password_reset.Destination) *ResetDestinationLog {
	// TODO あとでちゃんとする
	return &ResetDestinationLog{}
}
