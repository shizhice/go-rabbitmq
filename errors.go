package go_rabbitmq

import "errors"

var (
	ErrPoolClosed = errors.New("pool closed")
	ErrJobQueueIsNil = errors.New("job`s queue is nil")
)
