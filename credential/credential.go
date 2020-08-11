package credential

import (
	"time"

	"github.com/getto-systems/project-example-id/z_external/errors"
	"github.com/getto-systems/project-example-id/z_external/expiration"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	ErrClearCredential = errors.NewCategory("ClearCredential")

	ErrParseTicketMatchFailedNonce = NewClearCredentialError("Credential.ParseTicket", "MatchFailed.Nonce")
)

type (
	Credential struct {
		ticket       TicketToken
		apiToken     ApiToken
		contentToken ContentToken
	}

	TicketExtendSecond expiration.ExtendSecond

	TicketExpireSecond expiration.ExpireSecond
	TokenExpireSecond  expiration.ExpireSecond

	TicketExtendLimit expiration.ExtendLimit

	TicketExpires expiration.Expires
	TokenExpires  expiration.Expires

	Ticket struct {
		user  user.User
		nonce TicketNonce

		ticketExpires TicketExpires
		tokenExpires  TokenExpires
	}

	TicketSignature []byte
	TicketNonce     string
	TicketToken     struct {
		expires   TicketExpires
		nonce     TicketNonce
		signature TicketSignature
	}

	ApiRoles     []string
	ApiSignature []byte
	ApiToken     struct {
		expires   TokenExpires
		roles     ApiRoles
		signature ApiSignature
	}

	ContentKeyID     string
	ContentPolicy    string
	ContentSignature string
	ContentToken     struct {
		expires   TokenExpires
		keyID     ContentKeyID
		policy    ContentPolicy
		signature ContentSignature
	}
)

func NewCredential(ticket TicketToken, apiToken ApiToken, contentToken ContentToken) Credential {
	return Credential{
		ticket:       ticket,
		apiToken:     apiToken,
		contentToken: contentToken,
	}
}
func (credential Credential) TicketToken() TicketToken {
	return credential.ticket
}
func (credential Credential) ApiToken() ApiToken {
	return credential.apiToken
}
func (credential Credential) ContentToken() ContentToken {
	return credential.contentToken
}

func (expires TicketExpires) Expired(request request.Request) bool {
	return expiration.Expires(expires).Expired(time.Time(request.RequestedAt()))
}

func NewTicketExtendLimit(request request.Request, second TicketExtendSecond) TicketExtendLimit {
	return TicketExtendLimit(expiration.NewExtendLimit(
		time.Time(request.RequestedAt()),
		expiration.ExtendSecond(second),
	))
}
func (limit TicketExtendLimit) Extend(request request.Request, ticketSecond TicketExpireSecond, tokenSecond TokenExpireSecond) (TicketExpires, TokenExpires) {
	ticket := expiration.ExtendLimit(limit).Extend(time.Time(request.RequestedAt()), expiration.ExpireSecond(ticketSecond))
	token := expiration.ExtendLimit(limit).Extend(time.Time(request.RequestedAt()), expiration.ExpireSecond(tokenSecond))
	return TicketExpires(ticket), TokenExpires(token)
}

func NewTicketExpires(request request.Request, second TicketExpireSecond) TicketExpires {
	return TicketExpires(expiration.NewExpires(
		time.Time(request.RequestedAt()),
		expiration.ExpireSecond(second),
	))
}
func NewTokenExpires(request request.Request, second TokenExpireSecond) TokenExpires {
	return TokenExpires(expiration.NewExpires(
		time.Time(request.RequestedAt()),
		expiration.ExpireSecond(second),
	))
}

func NewTicket(user user.User, nonce TicketNonce, ticketExpires TicketExpires, tokenExpires TokenExpires) Ticket {
	return Ticket{
		user:  user,
		nonce: nonce,

		ticketExpires: ticketExpires,
		tokenExpires:  tokenExpires,
	}
}
func (ticket Ticket) User() user.User {
	return ticket.user
}
func (ticket Ticket) Nonce() TicketNonce {
	return ticket.nonce
}
func (ticket Ticket) TicketExpires() TicketExpires {
	return ticket.ticketExpires
}
func (ticket Ticket) TokenExpires() TokenExpires {
	return ticket.tokenExpires
}

func (ticket Ticket) NewTicketToken(signature TicketSignature) TicketToken {
	return TicketToken{
		expires:   ticket.ticketExpires,
		nonce:     ticket.nonce,
		signature: signature,
	}
}
func (ticket TicketToken) Expires() TicketExpires {
	return ticket.expires
}
func (ticket TicketToken) Nonce() TicketNonce {
	return ticket.nonce
}
func (ticket TicketToken) Signature() TicketSignature {
	return ticket.signature
}

func (ticket Ticket) NewApiToken(roles ApiRoles, signature ApiSignature) ApiToken {
	return ApiToken{
		expires:   ticket.tokenExpires,
		roles:     roles,
		signature: signature,
	}
}
func (token ApiToken) Expires() TokenExpires {
	return token.expires
}
func (token ApiToken) ApiRoles() ApiRoles {
	return token.roles
}
func (token ApiToken) Signature() ApiSignature {
	return token.signature
}

func (ticket Ticket) NewContentToken(keyID ContentKeyID, policy ContentPolicy, signature ContentSignature) ContentToken {
	return ContentToken{
		expires:   ticket.tokenExpires,
		keyID:     keyID,
		policy:    policy,
		signature: signature,
	}
}
func (token ContentToken) Expires() TokenExpires {
	return token.expires
}
func (token ContentToken) KeyID() ContentKeyID {
	return token.keyID
}
func (token ContentToken) Policy() ContentPolicy {
	return token.policy
}
func (token ContentToken) Signature() ContentSignature {
	return token.signature
}

func EmptyTicketExtendLimit() TicketExtendLimit {
	return TicketExtendLimit(expiration.EmptyExtendLimit())
}
func EmptyTicketExpires() TicketExpires {
	return TicketExpires(expiration.EmptyExpires())
}
func EmptyApiRoles() (_ ApiRoles) {
	return
}

func TicketExpireMinute(minutes int64) TicketExpireSecond {
	return TicketExpireSecond(minutes * 60)
}
func TicketExpireHour(hours int64) TicketExpireSecond {
	return TicketExpireMinute(hours * 60)
}
func TicketExpireDay(days int64) TicketExpireSecond {
	return TicketExpireHour(days * 24)
}
func TicketExpireWeek(days int64) TicketExpireSecond {
	return TicketExpireDay(days * 7)
}

func TicketExtendMinute(minutes int64) TicketExtendSecond {
	return TicketExtendSecond(minutes * 60)
}
func TicketExtendHour(hours int64) TicketExtendSecond {
	return TicketExtendMinute(hours * 60)
}
func TicketExtendDay(days int64) TicketExtendSecond {
	return TicketExtendHour(days * 24)
}
func TicketExtendWeek(days int64) TicketExtendSecond {
	return TicketExtendDay(days * 7)
}

func TokenExpireMinute(minutes int64) TokenExpireSecond {
	// トークンの有効期限は「分」の単位で設定するべき
	return TokenExpireSecond(minutes * 60)
}

func NewClearCredentialError(action string, message string) error {
	return errors.NewErrorAsCategory(action, message, ErrClearCredential)
}
