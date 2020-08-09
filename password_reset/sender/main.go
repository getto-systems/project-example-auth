package password_reset_sender

import (
	"errors"
	"fmt"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/message"
	"github.com/getto-systems/project-example-id/password_reset"
)

type (
	TokenSender struct {
		log message.LogMessage
	}
)

func NewTokenSender(log message.LogMessage) TokenSender {
	return TokenSender{
		log: log,
	}
}

func (sender TokenSender) tokenSender() infra.TokenSender {
	return sender
}

func (sender TokenSender) SendToken(dest password_reset.Destination, token password_reset.Token) (err error) {
	switch dest.Type {
	case "Log":
		return sender.log.Send(fmt.Sprintf("password reset token: %s", token))
	}

	return errors.New("unknown destination")
}
