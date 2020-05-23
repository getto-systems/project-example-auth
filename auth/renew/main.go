package renew

import (
	"time"

	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type Authenticator interface {
	UserRepository() user.UserRepository
	TicketTokener() token.TicketTokener
	AwsCloudFrontTokener() token.AwsCloudFrontTokener
}

type RenewParam struct {
	TicketToken token.TicketToken
	Path        user.Path
	Now         time.Time
}

func Renew(authenticator Authenticator, param RenewParam, handler auth.TokenHandler) (token.TicketInfo, error) {
	ticket, err := parseTicketToken(authenticator, param.TicketToken, param.Path)
	if err != nil {
		return nil, auth.ErrTicketTokenParseFailed
	}

	if ticket.IsRenewRequired(param.Now) {
		user, err := newUser(authenticator, ticket.User().UserID(), param.Path)
		if err != nil {
			return nil, auth.ErrUserAccessDenied
		}

		ticket = user.NewTicket(param.Now)

		ticketToken, err := ticketToken(authenticator, ticket)
		if err != nil {
			return nil, auth.ErrTicketTokenEncodeFailed
		}

		awsCloudFrontToken, err := awsCloudFrontToken(authenticator, ticket)
		if err != nil {
			return nil, auth.ErrAwsCloudFrontTokenEncodeFailed
		}

		handler(ticket, auth.Token{
			TicketToken:        ticketToken,
			AwsCloudFrontToken: awsCloudFrontToken,
		})
	}

	info, err := info(authenticator, ticket)
	if err != nil {
		return nil, auth.ErrInfoEncodeFailed
	}

	return info, nil
}

func parseTicketToken(authenticator Authenticator, ticketToken token.TicketToken, path user.Path) (user.Ticket, error) {
	return authenticator.TicketTokener().Parse(ticketToken, path)
}

func newUser(authenticator Authenticator, userID user.UserID, path user.Path) (user.User, error) {
	return user.NewUser(
		userID,
		authenticator.UserRepository().UserRoles(userID),
		path,
	)
}

func ticketToken(authenticator Authenticator, ticket user.Ticket) (token.TicketToken, error) {
	return authenticator.TicketTokener().Token(ticket)
}

func awsCloudFrontToken(authenticator Authenticator, ticket user.Ticket) (token.AwsCloudFrontToken, error) {
	return authenticator.AwsCloudFrontTokener().Token(ticket)
}

func info(authenticator Authenticator, ticket user.Ticket) (token.TicketInfo, error) {
	return authenticator.TicketTokener().Info(ticket)
}
