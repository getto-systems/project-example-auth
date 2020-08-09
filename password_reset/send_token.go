package password_reset

import (
	infra "github.com/getto-systems/project-example-id/infra/password_reset"
)

type SendToken struct {
	logger   infra.SendTokenLogger
	sessions infra.SessionRepository
	queue    infra.SendTokenJobQueue
	sender   infra.TokenSender
}

func NewSendToken(logger infra.SendTokenLogger, sessions infra.SessionRepository, queue infra.SendTokenJobQueue, sender infra.TokenSender) SendToken {
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
