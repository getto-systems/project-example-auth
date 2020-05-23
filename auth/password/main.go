package password

import (
	"time"

	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type Authenticator interface {
	UserRepository() user.UserRepository
	UserPasswordRepository() user.UserPasswordRepository
	TicketSerializer() token.TicketSerializer
	AwsCloudFrontSerializer() token.AwsCloudFrontSerializer
	Now() time.Time
}

type AuthParam struct {
	UserID       user.UserID
	UserPassword user.UserPassword
	Path         user.Path
}

func Auth(authenticator Authenticator, param AuthParam, handler auth.TokenHandler) (token.TicketInfo, error) {
	if !matchUserPassword(authenticator, param.UserID, param.UserPassword) {
		return nil, auth.ErrUserPasswordDidNotMatch
	}

	user, err := newUser(authenticator, param.UserID, param.Path)
	if err != nil {
		return nil, auth.ErrUserAccessDenied
	}

	ticket := user.NewTicket(authenticator.Now())

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

	info, err := ticketInfo(authenticator, ticket)
	if err != nil {
		return nil, auth.ErrTicketInfoSerializeFailed
	}

	return info, nil
}

func matchUserPassword(authenticator Authenticator, userID user.UserID, userPassword user.UserPassword) bool {
	return authenticator.UserPasswordRepository().MatchUserPassword(userID, userPassword)
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
