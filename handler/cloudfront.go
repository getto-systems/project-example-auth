package handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/auth"
)

func SetCloudFrontCookie(h AuthHandler, w http.ResponseWriter, ticket *auth.Ticket) error {
	token, err := h.CloudFrontSigner().Sign(ticket.Expires())
	if err != nil {
		return err
	}

	SetCookie(h, w, ticket, &Cookie{
		Name:  "CloudFront-Key-Pair-Id",
		Value: string(token.KeyPairID),
	})

	SetCookie(h, w, ticket, &Cookie{
		Name:  "CloudFront-Policy",
		Value: string(token.Policy),
	})

	SetCookie(h, w, ticket, &Cookie{
		Name:  "CloudFront-Signature",
		Value: string(token.Signature),
	})

	return nil
}
