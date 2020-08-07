package password_reset

import (
	"github.com/getto-systems/project-example-id/data/password_reset"
)

type SendToken struct {
	logger   password_reset.SendTokenLogger
	sessions password_reset.SessionRepository
	queue    password_reset.SendTokenJobQueue
	sender   password_reset.TokenSender
}

func NewSendToken(logger password_reset.SendTokenLogger, sessions password_reset.SessionRepository, queue password_reset.SendTokenJobQueue, sender password_reset.TokenSender) SendToken {
	return SendToken{
		logger:   logger,
		sessions: sessions,
		queue:    queue,
		sender:   sender,
	}
}

func (action SendToken) Send() (err error) {
	request, session, dest, token, err := action.queue.FetchSendTokenJob()
	if err != nil {
		return
	}

	action.logger.TryToSendToken(request, session, dest)

	err = action.sender.SendToken(dest, token)
	if err != nil {
		action.logger.FailedToSendToken(request, session, dest, err)

		updateErr := action.sessions.UpdateStatusToFailed(session, request.RequestedAt(), err)
		if updateErr != nil {
			// ここのステータスの更新失敗では送信エラーの内容を上書きしない
			action.logger.FailedToSendToken(request, session, dest, updateErr)
		}

		return
	}

	err = action.sessions.UpdateStatusToComplete(session, request.RequestedAt())
	if err != nil {
		action.logger.FailedToSendToken(request, session, dest, err)
		return
	}

	action.logger.SendToken(request, session, dest)
	return nil
}
