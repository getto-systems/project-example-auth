package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type ApiTokenIssuer struct {
	pub    EventPublisher
	signer ApiTokenSigner
	repo   userRolesRepository
}

func NewApiTokenIssuer(
	pub EventPublisher,
	db DB,
	signer ApiTokenSigner,
) ApiTokenIssuer {
	return ApiTokenIssuer{
		pub:    pub,
		signer: signer,
		repo: userRolesRepository{
			db: db,
		},
	}
}

func (issuer ApiTokenIssuer) issue(request data.Request, user data.User, expires data.Expires) (ApiToken, error) {
	roles := issuer.repo.roles(user)

	issuer.pub.IssueApiToken(request, user, roles, expires)

	token, err := issuer.signer.Sign(user, roles, expires)
	if err != nil {
		issuer.pub.IssueApiTokenFailed(request, user, roles, expires, err)
		return nil, err
	}

	return token, nil
}

type ApiTokenSigner interface {
	Sign(data.User, data.Roles, data.Expires) (ApiToken, error)
}

type userRolesRepository struct {
	db DB
}

func (repo userRolesRepository) roles(user data.User) data.Roles {
	roles, err := repo.db.FindUserRoles(user)
	if err != nil {
		// no permission when roles not found
		roles = nil
	}

	return roles
}
