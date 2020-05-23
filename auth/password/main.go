package password

import (
	"time"

	"github.com/getto-systems/project-example-id/auth"
	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
	"github.com/getto-systems/project-example-id/user/password"
)

type Authenticator interface {
	UserRepository() user.UserRepository
	UserPasswordRepository() password.UserPasswordRepository
	TicketTokener() token.TicketTokener
	AwsCloudFrontTokener() token.AwsCloudFrontTokener
}

type AuthParam struct {
	UserID       user.UserID
	UserPassword password.UserPassword
	Path         user.Path
	Now          time.Time
}

func Auth(authenticator Authenticator, param AuthParam, handler auth.TokenHandler) (token.TicketInfo, error) {
	if !matchUserPassword(authenticator, param.UserID, param.UserPassword) {
		return nil, auth.ErrUserPasswordDidNotMatch
	}

	user, err := newUser(authenticator, param.UserID, param.Path)
	if err != nil {
		return nil, auth.ErrUserAccessDenied
	}

	ticket := user.NewTicket(param.Now)

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

	info, err := info(authenticator, ticket)
	if err != nil {
		return nil, auth.ErrInfoEncodeFailed
	}

	return info, nil
}

func matchUserPassword(authenticator Authenticator, userID user.UserID, userPassword password.UserPassword) bool {
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
	return authenticator.TicketTokener().Token(ticket)
}

func awsCloudFrontToken(authenticator Authenticator, ticket user.Ticket) (token.AwsCloudFrontToken, error) {
	return authenticator.AwsCloudFrontTokener().Token(ticket)
}

func info(authenticator Authenticator, ticket user.Ticket) (token.TicketInfo, error) {
	return authenticator.TicketTokener().Info(ticket)
}
