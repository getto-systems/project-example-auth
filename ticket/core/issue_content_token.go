package core

import (
	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/ticket"
)

type contentTokenIssuer struct {
	pub    ticket.IssueContentTokenEventPublisher
	signer ticket.ContentTokenSigner
}

func newContentTokenIssuer(
	pub ticket.IssueContentTokenEventPublisher,
	signer ticket.ContentTokenSigner,
) contentTokenIssuer {
	return contentTokenIssuer{
		pub:    pub,
		signer: signer,
	}
}

func (issuer contentTokenIssuer) issue(request data.Request, user data.User, expires data.Expires) (ticket.ContentToken, error) {
	issuer.pub.IssueContentToken(request, user, expires)

	token, err := issuer.signer.Sign(expires)
	if err != nil {
		issuer.pub.IssueContentTokenFailed(request, user, expires, err)
		return nil, err
	}

	return token, nil
}
