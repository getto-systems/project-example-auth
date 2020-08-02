package api_token

type (
	ApiRoles     []string
	ApiSignature []byte
	ApiToken     struct {
		roles     ApiRoles
		signature ApiSignature
	}

	ContentKeyID     string
	ContentPolicy    string
	ContentSignature string
	ContentToken     struct {
		keyID     ContentKeyID
		policy    ContentPolicy
		signature ContentSignature
	}
)

func EmptyApiRoles() ApiRoles {
	return []string{}
}

func NewApiToken(roles ApiRoles, signature ApiSignature) ApiToken {
	return ApiToken{
		roles:     roles,
		signature: signature,
	}
}
func (token ApiToken) ApiRoles() ApiRoles {
	return token.roles
}
func (token ApiToken) Signature() ApiSignature {
	return token.signature
}

func NewContentToken(keyID ContentKeyID, policy ContentPolicy, signature ContentSignature) ContentToken {
	return ContentToken{
		keyID:     keyID,
		policy:    policy,
		signature: signature,
	}
}
func (token ContentToken) KeyID() ContentKeyID {
	return token.keyID
}
func (token ContentToken) Policy() ContentPolicy {
	return token.policy
}
func (token ContentToken) Signature() ContentSignature {
	return token.signature
}
