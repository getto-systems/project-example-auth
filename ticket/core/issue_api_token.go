package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type apiTokenIssuer struct {
	pub    ticket.IssueApiTokenEventPublisher
	signer ticket.ApiTokenSigner
	repo   userRolesRepository
}

func newApiTokenIssuer(
	pub ticket.IssueApiTokenEventPublisher,
	db ticket.IssueApiTokenDB,
	signer ticket.ApiTokenSigner,
) apiTokenIssuer {
	return apiTokenIssuer{
		pub:    pub,
		signer: signer,
		repo: userRolesRepository{
			db: db,
		},
	}
}

func (issuer apiTokenIssuer) issue(request data.Request, user data.User, expires data.Expires) (ticket.ApiToken, error) {
	roles := issuer.repo.roles(user)

	issuer.pub.IssueApiToken(request, user, roles, expires)

	token, err := issuer.signer.Sign(user, roles, expires)
	if err != nil {
		issuer.pub.IssueApiTokenFailed(request, user, roles, expires, err)
		return nil, err
	}

	return token, nil
}

type userRolesRepository struct {
	db ticket.IssueApiTokenDB
}

func (repo userRolesRepository) roles(user data.User) data.Roles {
	roles, err := repo.db.FindUserRoles(user)
	if err != nil {
		// no permission when roles not found
		roles = nil
	}

	return roles
}
