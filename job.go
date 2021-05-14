package go_rabbitmq

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"time"
)

type IJob interface {
	Config() *Config
	Queue() *Queue
	Exchange() *Exchange
	Handle(Payload)
	RoutingKey() string
}

// Dispatch message to queue
func Dispatch(job IJob, data Payload) error {
	if job.Queue() == nil {
		return ErrJobQueueIsNil
	}
	return publish(job, makeMessage(data, job), 0)
}

// Dispatch message to queue At time
func makeMessage(data Payload, job IJob) message {
	return message{
		Id: messageId(uuid.New().String()),
		Address: job.Config().getAddr().md5(),
		Queue: job.Queue().Name,
		Payload: data,
		State: messageSend,
		CreateAt: time.Now(),
	}
}


// Publish message to rabbitmq
func publish(job IJob, msg message, d TTL) error {
	// 创建MQ连接
	mq, err := Open(*job.Config())
	if err != nil {
		return err
	}

	activeQueue := job.Queue()
	if _, ok := job.(IDelayJob); ok {
		if activeQueue.Arguments == nil {
			activeQueue.Arguments = make(amqp.Table)
		}
		if job.Exchange() != nil {
			activeQueue.Arguments["x-dead-letter-exchange"] = job.Exchange()
		}
	}

	// 声明队列
	if err := declareJob(job, mq); err != nil {
		return err
	}

	return (&postman{
		mq: mq,
		job: job,
		msg: msg,
		ttl: d,
	}).send()
}

func declareJob(job IJob, mq *MQ) error {
	list, err := shouldDeclareQueue(job)
	if err != nil {
		return err
	}

	for i := range list {
		if err = declareQueue(list[i], mq); err != nil {
			return err
		}
	}

	return nil
}

func declareQueue(queue *Queue, mq *MQ) error {
	if queue.Binding.Exchange == nil {
		queue.Binding.Exchange = commonTopicExchange
	}
	// 声明交换机
	if err := mq.DeclareExchange(queue.Binding.Exchange); err != nil {
		return err
	}

	// 声明队列
	if err := mq.DeclareQueue(queue); err != nil {
		return err
	}

	return nil
}

func shouldDeclareQueue(job IJob) (list []*Queue, err error) {
	activeQueue := job.Queue()
	list = append(list, activeQueue)
	if job, ok := job.(IDelayJob); ok {
		var dlxExchange = deadLetterExchange
		var dlxQueue *Queue
		defer func() {
			if dlxQueue != nil {
				list = append(list, dlxQueue)
			}
		}()
		if activeQueue.Arguments == nil {
			activeQueue.Arguments = make(amqp.Table)
		}

		if job.DlxExchange() != nil {
			dlxExchange = job.DlxExchange()
		}
		activeQueue.Arguments["x-dead-letter-exchange"] = dlxExchange.Name

		if job.DlxQueue() != nil {
			dlxQueue = job.DlxQueue()
			activeQueue.Arguments["x-dead-letter-queue"] = dlxQueue.Name
		} else {
			var args = make(amqp.Table)
			if activeQueue.Arguments != nil {
				for key := range activeQueue.Arguments {
					if key != "x-dead-letter-exchange" && key != "x-dead-letter-queue" && key != "x-message-ttl" {
						args[key] = activeQueue.Arguments[key]
					}
				}
			}

			dlxQueue = &Queue{
				Name: fmt.Sprintf("%s.dlx", activeQueue.Name),
				Binding: QueueBind{
					Exchange: dlxExchange,
					BindingKey: activeQueue.Binding.BindingKey,
					Args: activeQueue.Binding.Args,
				},
				Arguments: args,
			}
		}

		activeQueue.Arguments["x-dead-letter-queue"] = dlxQueue.Name

		if job.TTL() >= TTLMillisecond {
			activeQueue.Arguments["x-message-ttl"] = job.TTL().Value()
		}
	}

	return list, nil
}