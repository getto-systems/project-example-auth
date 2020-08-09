package client

import (
	credential_infra "github.com/getto-systems/project-example-id/infra/credential"
	password_infra "github.com/getto-systems/project-example-id/infra/password"
	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"
	ticket_infra "github.com/getto-systems/project-example-id/infra/ticket"
	user_infra "github.com/getto-systems/project-example-id/infra/user"

	password_reset_data "github.com/getto-systems/project-example-id/data/password_reset"
	ticket_data "github.com/getto-systems/project-example-id/data/ticket"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/ticket"
	"github.com/getto-systems/project-example-id/user"
)

type (
	Backend struct {
		ticket        TicketAction
		credential    CredentialAction
		user          UserAction
		password      PasswordAction
		passwordReset PasswordResetAction
	}

	TicketAction struct {
		register   ticket.Register
		validate   ticket.Validate
		extend     ticket.Extend
		deactivate ticket.Deactivate
	}

	CredentialAction struct {
		parseTicket       credential.ParseTicket
		issueTicket       credential.IssueTicket
		issueApiToken     credential.IssueApiToken
		issueContentToken credential.IssueContentToken
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
	credential CredentialAction,
	user UserAction,
	password PasswordAction,
	passwordReset PasswordResetAction,
) Backend {
	return Backend{
		ticket:        ticket,
		credential:    credential,
		user:          user,
		password:      password,
		passwordReset: passwordReset,
	}
}

func NewTicketAction(
	logger ticket_infra.Logger,

	gen credential_infra.TicketNonceGenerator,

	tickets ticket_infra.TicketRepository,
) TicketAction {
	return TicketAction{
		register:   ticket.NewRegister(logger, gen, tickets),
		validate:   ticket.NewValidate(logger, tickets),
		extend:     ticket.NewExtend(logger, tickets),
		deactivate: ticket.NewDeactivate(logger, tickets),
	}
}

func NewCredentialAction(
	logger credential_infra.Logger,

	ticketSign credential_infra.TicketSign,
	apiTokenSinger credential_infra.ApiTokenSigner,
	contentTokenSigner credential_infra.ContentTokenSigner,

	apiUsers credential_infra.ApiUserRepository,
) CredentialAction {
	return CredentialAction{
		parseTicket:       credential.NewParseTicket(logger, ticketSign),
		issueTicket:       credential.NewIssueTicket(logger, ticketSign),
		issueApiToken:     credential.NewIssueApiToken(logger, apiTokenSinger, apiUsers),
		issueContentToken: credential.NewIssueContentToken(logger, contentTokenSigner),
	}
}

func NewUserAction(
	logger user_infra.Logger,

	users user_infra.UserRepository,
) UserAction {
	return UserAction{
		getLogin: user.NewGetLogin(logger, users),
		getUser:  user.NewGetUser(logger, users),
	}
}

func NewPasswordAction(
	logger password_infra.Logger,

	exp ticket_data.Expiration,
	enc password_infra.PasswordEncrypter,

	passwords password_infra.PasswordRepository,
) PasswordAction {
	return PasswordAction{
		validate: password.NewValidate(logger, exp, enc, passwords),
		change:   password.NewChange(logger, enc, passwords),
	}
}

func NewPasswordResetAction(
	logger password_reset_infra.Logger,

	ticketExp ticket_data.Expiration,
	exp password_reset_data.Expiration,
	gen password_reset_infra.SessionGenerator,

	sessions password_reset_infra.SessionRepository,
	destinations password_reset_infra.DestinationRepository,

	tokenQueue password_reset_infra.SendTokenJobQueue,

	tokenSender password_reset_infra.TokenSender,
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
