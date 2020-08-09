package client

import (
	credential_infra "github.com/getto-systems/project-example-id/infra/credential"
	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/credential"
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/user"

	credential_core "github.com/getto-systems/project-example-id/credential"
	password_reset_core "github.com/getto-systems/project-example-id/password_reset"
)

type (
	Backend struct {
		ticket        ticket.Action
		credential    CredentialAction
		user          user.Action
		password      password.Action
		passwordReset PasswordResetAction
	}

	CredentialAction struct {
		parseTicket       credential.ParseTicket
		issueTicket       credential.IssueTicket
		issueApiToken     credential.IssueApiToken
		issueContentToken credential.IssueContentToken
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
	ticket ticket.Action,
	credential CredentialAction,
	user user.Action,
	password password.Action,
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
