package renew

import (
	"time"

	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type Authenticator interface {
	UserRepository() user.UserRepository
	TicketSerializer() token.TicketSerializer
	AwsCloudFrontSerializer() token.AwsCloudFrontSerializer
	Now() time.Time
}

type RenewParam struct {
	TicketToken token.TicketToken
	Path        user.Path
}

func Renew(authenticator Authenticator, param RenewParam, handler auth.TokenHandler) (token.TicketInfo, error) {
	ticket, err := parseTicketToken(authenticator, param.TicketToken, param.Path)
	if err != nil {
		return nil, auth.ErrTicketTokenParseFailed
	}

	now := authenticator.Now()

	if ticket.IsRenewRequired(now) {
		user, err := newUser(authenticator, ticket.User().UserID(), param.Path)
		if err != nil {
			return nil, auth.ErrUserAccessDenied
		}

		ticket = user.NewTicket(now)

		ticketToken, err := ticketToken(authenticator, ticket)
		if err != nil {
			return nil, auth.ErrTicketTokenSerializeFailed
		}

		awsCloudFrontToken, err := awsCloudFrontToken(authenticator, ticket)
		if err != nil {
			return nil, auth.ErrAwsCloudFrontTokenSerializeFailed
		}

		handler(ticket, auth.Token{
			TicketToken:        ticketToken,
			AwsCloudFrontToken: awsCloudFrontToken,
		})
	}

	info, err := ticketInfo(authenticator, ticket)
	if err != nil {
		return nil, auth.ErrTicketInfoSerializeFailed
	}

	return info, nil
}

func parseTicketToken(authenticator Authenticator, ticketToken token.TicketToken, path user.Path) (user.Ticket, error) {
	return authenticator.TicketSerializer().Parse(ticketToken, path)
}

func newUser(authenticator Authenticator, userID user.UserID, path user.Path) (user.User, error) {
	return user.NewUser(
		userID,
		authenticator.UserRepository().UserRoles(userID),
		path,
	)
}

func ticketToken(authenticator Authenticator, ticket user.Ticket) (token.TicketToken, error) {
	return authenticator.TicketSerializer().Token(ticket)
}

func awsCloudFrontToken(authenticator Authenticator, ticket user.Ticket) (token.AwsCloudFrontToken, error) {
	return authenticator.AwsCloudFrontSerializer().Token(ticket)
}

func ticketInfo(authenticator Authenticator, ticket user.Ticket) (token.TicketInfo, error) {
	return authenticator.TicketSerializer().Info(ticket)
}
