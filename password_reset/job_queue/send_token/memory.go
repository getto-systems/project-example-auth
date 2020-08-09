package password_reset_job_queue_send_token

import (
	"errors"

	"github.com/getto-systems/project-example-id/password_reset/infra"

	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/password_reset"
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

func (queue *MemoryQueue) queue() infra.SendTokenJobQueue {
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
