package go_rabbitmq

import (
	"github.com/streadway/amqp"
)


type postman struct {
	mq *MQ
	job IJob
	msg message
	ttl TTL
	notifyConfirmChan chan amqp.Confirmation
	notifyReturnChan  chan amqp.Return
}

func (p *postman) send() error {
	activeCh, err := p.mq.pool.get()
	if err != nil {
		return err
	}
	defer func() {
		p.mq.pool.release(activeCh)
	}()

	exchangeName := commonTopicExchange.Name
	if p.job.Exchange() != nil {
		exchangeName = p.job.Exchange().Name
	}

	return activeCh.publish(exchangeName, p.job.RoutingKey(), p.msg, p.ttl)
}
