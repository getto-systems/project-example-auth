package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type ContentTokenIssuer struct {
	pub    EventPublisher
	signer ContentTokenSigner
}

func NewContentTokenIssuer(
	pub EventPublisher,
	signer ContentTokenSigner,
) ContentTokenIssuer {
	return ContentTokenIssuer{
		pub:    pub,
		signer: signer,
	}
}

func (issuer ContentTokenIssuer) issue(request data.Request, user data.User, expires data.Expires) (ContentToken, error) {
	issuer.pub.IssueContentToken(request, user, expires)

	token, err := issuer.signer.Sign(expires)
	if err != nil {
		issuer.pub.IssueContentTokenFailed(request, user, expires, err)
		return nil, err
	}

	return token, nil
}

type ContentTokenSigner interface {
	Sign(data.Expires) (ContentToken, error)
}
