package static_repository_secret

import (
	"context"
	"encoding/json"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/getto-systems/project-example-auth/y_static/infra"

	"github.com/getto-systems/project-example-auth/y_static"
)

type (
	SecretManager struct {
	}
)

func NewSecretManager() SecretManager {
	return SecretManager{}
}

func (repo SecretManager) repo() infra.SecretRepository {
	return repo
}

func (SecretManager) FindSecret(name string) (_ static.Secret, err error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return
	}

	type (
		AdminSecret struct {
			UserID   string `json:"user_id"`
			LoginID  string `json:"login_id"`
			Password string `json:"password"`
		}
		CookieSecret struct {
			Domain string `json:"domain"`
		}
		TicketSecret struct {
			PrivateKey string `json:"private_key"`
			PublicKey  string `json:"public_key"`
		}
		ApiSecret struct {
			PrivateKey string `json:"private_key"`
		}
		CloudfrontSecret struct {
			KeyPairID   string `json:"key_pair_id"`
			PrivateKey  string `json:"private_key"`
			ResourceURL string `json:"resource_url"`
		}
		Secret struct {
			Admin      AdminSecret      `json:"admin"`
			Cookie     CookieSecret     `json:"cookie"`
			Ticket     TicketSecret     `json:"ticket"`
			Api        ApiSecret        `json:"api"`
			Cloudfront CloudfrontSecret `json:"cloudfront"`
		}
	)

	var input Secret
	err = json.Unmarshal(result.Payload.Data, &input)
	if err != nil {
		return
	}

	return static.Secret{
		Admin: static.AdminSecret{
			UserID:   input.Admin.UserID,
			LoginID:  input.Admin.LoginID,
			Password: input.Admin.Password,
		},
		Cookie: static.CookieSecret{
			Domain: input.Cookie.Domain,
		},
		Ticket: static.TicketSecret{
			PrivateKey: []byte(input.Ticket.PrivateKey),
			PublicKey:  []byte(input.Ticket.PublicKey),
		},
		Api: static.ApiSecret{
			PrivateKey: []byte(input.Api.PrivateKey),
		},
		Cloudfront: static.CloudfrontSecret{
			KeyPairID:   input.Cloudfront.KeyPairID,
			PrivateKey:  []byte(input.Cloudfront.PrivateKey),
			ResourceURL: input.Cloudfront.ResourceURL,
		},
	}, nil
}
