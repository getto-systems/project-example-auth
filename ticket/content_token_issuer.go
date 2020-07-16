package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type ContentTokenIssuer struct {
	pub    contentTokenIssueEventPublisher
	signer ContentTokenSigner
}

type contentTokenIssueEventPublisher interface {
	IssueContentToken(data.Request, data.User, data.Expires)
	IssueContentTokenFailed(data.Request, data.User, data.Expires, error)
}

func NewContentTokenIssuer(
	pub contentTokenIssueEventPublisher,
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
