package go_rabbitmq

import (
	"fmt"
	"time"
)

type TTL time.Duration

func (t TTL) String() string {
	return fmt.Sprintf("%d", t / 1000 / 1000)
}

func (t TTL) Value() int {
	return int(t / 1000 / 1000)
}

var (
	TTLMillisecond = 1000 * TTL(time.Microsecond)
	TTLSecond      = 1000 * TTLMillisecond
	TTLMinute      = 60 * TTLSecond
	TTLHour        = 60 * TTLMinute
	TTLDay         = 24 * TTLHour
	TTLWeek        = 7 * TTLDay
	MixTTL         = TTLMillisecond
)

type IDelayJob interface {
	IJob
	TTL() TTL
	DlxExchange() *Exchange
	DlxQueue() *Queue
}

// delay Dispatch message to queue
func DispatchDelay(job IDelayJob, data Payload) error {
	if job.Queue() == nil {
		return ErrJobQueueIsNil
	}
	return publish(job, makeMessage(data, job), 0)
}

// Dispatch message to queue At time
func delayTime(job IDelayJob) TTL {
	return job.TTL()
}
