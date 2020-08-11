package logger

import (
	"encoding/json"
	golog "log"
	"os"
	"time"

	"github.com/getto-systems/applog-go/v2"

	"github.com/getto-systems/project-example-id/_gateway/log"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Logger struct {
		logger applog.Logger
	}
)

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

type (
	Entry struct {
		Time    string     `json:"time"`
		Level   string     `json:"level"`
		Message string     `json:"message"`
		Request RequestLog `json:"request"`

		User  *UserLog  `json:"user,omitempty"`
		Login *LoginLog `json:"login,omitempty"`

		Credential    *CredentialEntry    `json:"credential,omitempty"`
		PasswordReset *PasswordResetEntry `json:"password_reset,omitempty"`

		Error string `json:"error,omitempty"`
	}

	CredentialEntry struct {
		ApiRoles *ApiRolesLog `json:"roles,omitempty"`

		TicketExpires     *CredentialTicketExpiresLog     `json:"ticket_expires,omitempty"`
		TicketExtendLimit *CredentialTicketExtendLimitLog `json:"ticket_limit,omitempty"`

		TokenExpires *CredentialTokenExpiresLog `json:"token_expires,omitempty"`
	}

	PasswordResetEntry struct {
		Session     *ResetSessionLog        `json:"session,omitempty"`
		Status      *ResetStatusLog         `json:"status,omitempty"`
		Destination *ResetDestinationLog    `json:"destination,omitempty"`
		Expires     *ResetSessionExpiresLog `json:"expires,omitempty"`
	}

	RequestLog struct {
		RequestedAt string   `json:"requested_at"`
		Route       RouteLog `json:"route"`
	}
	RouteLog struct {
		RemoteAddr string `json:"remote_addr"`
	}

	UserLog struct {
		UserID string `json:"user_id"`
	}
	LoginLog struct {
		LoginID string `json:"login_id"`
	}

	ApiRolesLog struct {
		ApiRoles []string `json:"api_roles"`
	}

	CredentialTicketExpiresLog struct {
		Expires string `json:"expires"`
	}
	CredentialTicketExtendLimitLog struct {
		ExtendLimit string `json:"extend_limit"`
	}

	CredentialTokenExpiresLog struct {
		Expires string `json:"expires"`
	}

	ResetSessionLog struct {
		SessionID string `json:"session_id"`
	}

	// TODO あとでちゃんとする
	ResetStatusLog struct {
	}

	// TODO あとでちゃんとする
	ResetDestinationLog struct {
	}

	ResetSessionExpiresLog struct {
		Expires string `json:"expires"`
	}
)

func format(log log.Entry) Entry {
	entry := Entry{
		Time:    time.Now().UTC().String(),
		Message: log.Message,
	}

	entry.Request = requestLog(log.Request)

	if log.User != nil {
		entry.User = userLog(log.User)
	}
	if log.Login != nil {
		entry.Login = loginLog(log.Login)
	}

	if log.Credential != nil {
		entry.Credential = credentialLog(log.Credential)
	}
	if log.PasswordReset != nil {
		entry.PasswordReset = passwordResetLog(log.PasswordReset)
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
		RequestedAt: time.Time(request.RequestedAt()).String(),
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

func credentialLog(log *log.CredentialEntry) (entry *CredentialEntry) {
	if log.ApiRoles != nil {
		entry.ApiRoles = apiRolesLog(log.ApiRoles)
	}

	if log.TicketExpires != nil {
		entry.TicketExpires = credentialTicketExpiresLog(log.TicketExpires)
	}
	if log.TicketExtendLimit != nil {
		entry.TicketExtendLimit = credentialTicketExtendLimitLog(log.TicketExtendLimit)
	}

	if log.TokenExpires != nil {
		entry.TokenExpires = credentialTokenExpiresLog(log.TokenExpires)
	}

	return
}
func apiRolesLog(roles *credential.ApiRoles) (log *ApiRolesLog) {
	for _, role := range *roles {
		log.ApiRoles = append(log.ApiRoles, string(role))
	}
	return
}
func credentialTicketExpiresLog(expires *credential.TicketExpires) *CredentialTicketExpiresLog {
	return &CredentialTicketExpiresLog{
		Expires: time.Time(*expires).String(),
	}
}
func credentialTicketExtendLimitLog(limit *credential.TicketExtendLimit) *CredentialTicketExtendLimitLog {
	return &CredentialTicketExtendLimitLog{
		ExtendLimit: time.Time(*limit).String(),
	}
}
func credentialTokenExpiresLog(expires *credential.TokenExpires) *CredentialTokenExpiresLog {
	return &CredentialTokenExpiresLog{
		Expires: time.Time(*expires).String(),
	}
}

func passwordResetLog(log *log.PasswordResetEntry) (entry *PasswordResetEntry) {
	if log.Session != nil {
		entry.Session = resetSessionLog(log.Session)
	}
	if log.Status != nil {
		entry.Status = resetStatusLog(log.Status)
	}
	if log.Destination != nil {
		entry.Destination = resetDestinationLog(log.Destination)
	}
	if log.Expires != nil {
		entry.Expires = resetSessionExpiresLog(log.Expires)
	}
	return
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
func resetSessionExpiresLog(expires *password_reset.Expires) *ResetSessionExpiresLog {
	return &ResetSessionExpiresLog{
		Expires: time.Time(*expires).String(),
	}
}
