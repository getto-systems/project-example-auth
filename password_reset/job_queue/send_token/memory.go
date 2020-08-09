package send_token

import (
	"errors"

	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
)

type (
	MemoryQueue struct {
		jobs []job
	}

	job struct {
		request     request.Request
		session     password_reset.Session
		destination password_reset.Destination
		token       password_reset.Token
	}
)

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{}
}

func (queue *MemoryQueue) queue() password_reset_infra.SendTokenJobQueue {
	return queue
}

func (queue *MemoryQueue) PushSendTokenJob(request request.Request, session password_reset.Session, dest password_reset.Destination, token password_reset.Token) (err error) {
	queue.jobs = append(queue.jobs, job{
		request:     request,
		session:     session,
		destination: dest,
		token:       token,
	})
	return nil
}

func (queue *MemoryQueue) FetchSendTokenJob() (_ request.Request, _ password_reset.Session, _ password_reset.Destination, _ password_reset.Token, err error) {
	if len(queue.jobs) == 0 {
		err = errors.New("empty queue")
		return
	}

	job := queue.jobs[0]
	queue.jobs = queue.jobs[1:]

	return job.request, job.session, job.destination, job.token, nil
}
