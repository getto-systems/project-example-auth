package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type (
	ResetStatus struct {
		requestedAt  *data.Time
		deliveringAt *data.Time
		delivered    *ResetDelivered
		deliverError *ResetDeliverError
	}
	ResetDelivered struct {
		deliveredAt data.Time
		destination string
	}
	ResetDeliverError struct {
		erroredAt data.Time
		err       string
	}
)

func NewResetStatusNotRequested() ResetStatus {
	return ResetStatus{}
}

func NewResetStatusRequested(requestedAt data.Time) ResetStatus {
	return ResetStatus{
		requestedAt: &requestedAt,
	}
}

func NewResetStatusDelivering(requestedAt data.Time, deliveringAt data.Time) ResetStatus {
	return ResetStatus{
		requestedAt:  &requestedAt,
		deliveringAt: &deliveringAt,
	}
}

func NewResetStatusDelivered(requestedAt data.Time, deliveringAt data.Time, deliveredAt data.Time, destination string) ResetStatus {
	return ResetStatus{
		requestedAt:  &requestedAt,
		deliveringAt: &deliveringAt,
		delivered:    newResetDelivered(deliveredAt, destination),
	}
}

func NewResetStatusDeliverError(requestedAt data.Time, deliveringAt data.Time, erroredAt data.Time, err error) ResetStatus {
	return ResetStatus{
		requestedAt:  &requestedAt,
		deliveringAt: &deliveringAt,
		deliverError: newResetDeliverError(erroredAt, err),
	}
}

func (status ResetStatus) RequestedAt() *data.Time {
	return status.requestedAt
}

func (status ResetStatus) Delivering() *data.Time {
	return status.deliveringAt
}

func (status ResetStatus) Delivered() *ResetDelivered {
	return status.delivered
}

func (status ResetStatus) DeliverError() *ResetDeliverError {
	return status.deliverError
}

func newResetDelivered(deliveredAt data.Time, destination string) *ResetDelivered {
	return &ResetDelivered{
		deliveredAt: deliveredAt,
		destination: destination,
	}
}

func (status ResetDelivered) DeliveredAt() data.Time {
	return status.deliveredAt
}

func (status ResetDelivered) Destination() string {
	return status.destination
}

func newResetDeliverError(erroredAt data.Time, err error) *ResetDeliverError {
	return &ResetDeliverError{
		erroredAt: erroredAt,
		err:       err.Error(),
	}
}

func (status ResetDeliverError) ErroredAt() data.Time {
	return status.erroredAt
}

func (status ResetDeliverError) Error() string {
	return status.err
}
