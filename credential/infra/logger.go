package infra

import (
	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/request"
	"github.com/getto-systems/project-example-auth/user"
)

type (
	Logger interface {
		ParseTicketSignatureLogger

		IssueTicketTokenLogger
		IssueApiTokenLogger
		IssueContentTokenLogger
	}

	ParseTicketSignatureLogger interface {
		TryToParseTicketSignature(request.Request)
		FailedToParseTicketSignature(request.Request, error)
		FailedToParseTicketSignatureBecauseNonceMatchFailed(request.Request, error)
		ParseTicketSignature(request.Request, user.User)
	}

	IssueTicketTokenLogger interface {
		TryToIssueTicketToken(request.Request, user.User, credential.TicketExpires)
		FailedToIssueTicketToken(request.Request, user.User, credential.TicketExpires, error)
		IssueTicketToken(request.Request, user.User, credential.TicketExpires)
	}

	IssueApiTokenLogger interface {
		TryToIssueApiToken(request.Request, user.User, credential.TokenExpires)
		FailedToIssueApiToken(request.Request, user.User, credential.TokenExpires, error)
		IssueApiToken(request.Request, user.User, credential.ApiRoles, credential.TokenExpires)
	}

	IssueContentTokenLogger interface {
		TryToIssueContentToken(request.Request, user.User, credential.TokenExpires)
		FailedToIssueContentToken(request.Request, user.User, credential.TokenExpires, error)
		IssueContentToken(request.Request, user.User, credential.TokenExpires)
	}
)
