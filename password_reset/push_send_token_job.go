package password_reset

import (
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

type PushSendTokenJob struct {
	logger    password_reset.PushSendTokenJobLogger
	sessions  password_reset.SessionRepository
	sendQueue password_reset.SendTokenJobQueue
}

func NewPushSendTokenJob(logger password_reset.PushSendTokenJobLogger, sessions password_reset.SessionRepository, sendQueue password_reset.SendTokenJobQueue) PushSendTokenJob {
	return PushSendTokenJob{
		logger:    logger,
		sessions:  sessions,
		sendQueue: sendQueue,
	}
}

func (action PushSendTokenJob) Push(request request.Request, session password_reset.Session, dest password_reset.Destination, token password_reset.Token) (err error) {
	action.logger.TryToPushSendTokenJob(request, session, dest)

	err = action.sessions.UpdateStatusToSending(session, request.RequestedAt())
	if err != nil {
		action.logger.FailedToPushSendTokenJob(request, session, dest, err)
		return
	}

	err = action.sendQueue.PushSendTokenJob(request, session, dest, token)
	if err != nil {
		action.logger.FailedToPushSendTokenJob(request, session, dest, err)

		updateErr := action.sessions.UpdateStatusToFailed(session, request.RequestedAt(), err)
		if updateErr != nil {
			// ここのステータスの更新失敗では送信エラーの内容を上書きしない
			action.logger.FailedToPushSendTokenJob(request, session, dest, updateErr)
		}

		return
	}

	// 送信済みにするのは worker の仕事
	action.logger.PushSendTokenJob(request, session, dest)
	return nil
}
