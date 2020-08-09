package client

import (
	credential_infra "github.com/getto-systems/project-example-id/infra/credential"
	password_infra "github.com/getto-systems/project-example-id/infra/password"
	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"
	ticket_infra "github.com/getto-systems/project-example-id/infra/ticket"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"

	credential_core "github.com/getto-systems/project-example-id/credential"
	password_core "github.com/getto-systems/project-example-id/password"
	password_reset_core "github.com/getto-systems/project-example-id/password_reset"
	ticket_core "github.com/getto-systems/project-example-id/ticket"
)

type (
	Backend struct {
		ticket        TicketAction
		credential    CredentialAction
		user          user.Action
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
	user user.Action,
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
		register:   ticket_core.NewRegister(logger, gen, tickets),
		validate:   ticket_core.NewValidate(logger, tickets),
		extend:     ticket_core.NewExtend(logger, tickets),
		deactivate: ticket_core.NewDeactivate(logger, tickets),
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
		parseTicket:       credential_core.NewParseTicket(logger, ticketSign),
		issueTicket:       credential_core.NewIssueTicket(logger, ticketSign),
		issueApiToken:     credential_core.NewIssueApiToken(logger, apiTokenSinger, apiUsers),
		issueContentToken: credential_core.NewIssueContentToken(logger, contentTokenSigner),
	}
}

func NewPasswordAction(
	logger password_infra.Logger,

	exp ticket.Expiration,
	enc password_infra.PasswordEncrypter,

	passwords password_infra.PasswordRepository,
) PasswordAction {
	return PasswordAction{
		validate: password_core.NewValidate(logger, exp, enc, passwords),
		change:   password_core.NewChange(logger, enc, passwords),
	}
}

func NewPasswordResetAction(
	logger password_reset_infra.Logger,

	ticketExp ticket.Expiration,
	exp password_reset.Expiration,
	gen password_reset_infra.SessionGenerator,

	sessions password_reset_infra.SessionRepository,
	destinations password_reset_infra.DestinationRepository,

	tokenQueue password_reset_infra.SendTokenJobQueue,

	tokenSender password_reset_infra.TokenSender,
) PasswordResetAction {
	return PasswordResetAction{
		createSession:    password_reset_core.NewCreateSession(logger, exp, gen, sessions, destinations),
		pushSendTokenJob: password_reset_core.NewPushSendTokenJob(logger, sessions, tokenQueue),
		sendToken:        password_reset_core.NewSendToken(logger, sessions, tokenQueue, tokenSender),
		getStatus:        password_reset_core.NewGetStatus(logger, sessions),
		validate:         password_reset_core.NewValidate(logger, ticketExp, sessions),
		closeSession:     password_reset_core.NewCloseSession(logger, sessions),
	}
}
