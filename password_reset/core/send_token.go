package password_reset_core

func (action action) SendToken() (err error) {
	request, session, dest, token, err := action.tokenQueue.FetchSendTokenJob()
	if err != nil {
		return
	}

	action.logger.TryToSendToken(request, session, dest)

	err = action.tokenSender.SendToken(dest, token)
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
