package client

import (
	api_token_data "github.com/getto-systems/project-example-id/data/api_token"
	password_data "github.com/getto-systems/project-example-id/data/password"
	password_reset_data "github.com/getto-systems/project-example-id/data/password_reset"
	ticket_data "github.com/getto-systems/project-example-id/data/ticket"
	user_data "github.com/getto-systems/project-example-id/data/user"

	"github.com/getto-systems/project-example-id/api_token"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/ticket"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Backend struct {
		ticket        TicketAction
		apiToken      ApiTokenAction
		user          UserAction
		password      PasswordAction
		passwordReset PasswordResetAction
	}

	TicketAction struct {
		validate   ticket.Validate
		extend     ticket.Extend
		deactivate ticket.Deactivate
		issue      ticket.Issue
	}

	ApiTokenAction struct {
		issueApiToken     api_token.IssueApiToken
		issueContentToken api_token.IssueContentToken
	}

	UserAction struct {
		getLogin user.GetLogin
		getUser  user.GetUser
	}

	PasswordAction struct {
		validate password.Validate
		change   password.Change
	}

	PasswordResetAction struct {
		createSession    password_reset.CreateSession
		pushSendTokenJob password_reset.PushSendTokenJob
		sendToken        password_reset.SendToken
		getStatus        password_reset.GetStatus
		validate         password_reset.Validate
		closeSession     password_reset.CloseSession
	}
)

func NewBackend(
	ticket TicketAction,
	apiToken ApiTokenAction,
	user UserAction,
	password PasswordAction,
	passwordReset PasswordResetAction,
) Backend {
	return Backend{
		ticket:        ticket,
		apiToken:      apiToken,
		user:          user,
		password:      password,
		passwordReset: passwordReset,
	}
}

func NewTicketAction(
	logger ticket_data.Logger,

	sign api_token_data.TicketSign,
	gen api_token_data.TicketNonceGenerator,

	tickets ticket_data.TicketRepository,
) TicketAction {
	return TicketAction{
		validate:   ticket.NewValidate(logger, sign, tickets),
		extend:     ticket.NewExtend(logger, sign, tickets),
		deactivate: ticket.NewDeactivate(logger, tickets),
		issue:      ticket.NewIssue(logger, sign, gen, tickets),
	}
}

func NewApiTokenAction(
	logger api_token_data.Logger,

	apiTokenSinger api_token_data.ApiTokenSigner,
	contentTokenSigner api_token_data.ContentTokenSigner,

	apiUsers api_token_data.ApiUserRepository,
) ApiTokenAction {
	return ApiTokenAction{
		issueApiToken:     api_token.NewIssueApiToken(logger, apiTokenSinger, apiUsers),
		issueContentToken: api_token.NewIssueContentToken(logger, contentTokenSigner),
	}
}

func NewUserAction(
	logger user_data.Logger,

	users user_data.UserRepository,
) UserAction {
	return UserAction{
		getLogin: user.NewGetLogin(logger, users),
		getUser:  user.NewGetUser(logger, users),
	}
}

func NewPasswordAction(
	logger password_data.Logger,

	exp ticket_data.Expiration,
	enc password_data.PasswordEncrypter,

	passwords password_data.PasswordRepository,
) PasswordAction {
	return PasswordAction{
		validate: password.NewValidate(logger, exp, enc, passwords),
		change:   password.NewChange(logger, enc, passwords),
	}
}

func NewPasswordResetAction(
	logger password_reset_data.Logger,

	ticketExp ticket_data.Expiration,
	exp password_reset_data.Expiration,
	gen password_reset_data.SessionGenerator,

	sessions password_reset_data.SessionRepository,
	destinations password_reset_data.DestinationRepository,

	tokenQueue password_reset_data.SendTokenJobQueue,

	tokenSender password_reset_data.TokenSender,
) PasswordResetAction {
	return PasswordResetAction{
		createSession:    password_reset.NewCreateSession(logger, exp, gen, sessions, destinations),
		pushSendTokenJob: password_reset.NewPushSendTokenJob(logger, sessions, tokenQueue),
		sendToken:        password_reset.NewSendToken(logger, sessions, tokenQueue, tokenSender),
		getStatus:        password_reset.NewGetStatus(logger, sessions),
		validate:         password_reset.NewValidate(logger, ticketExp, sessions),
		closeSession:     password_reset.NewCloseSession(logger, sessions),
	}
}
