package go_rabbitmq

import (
	"github.com/streadway/amqp"
	"time"
)

type channel struct {
	ch *amqp.Channel
	createdAt time.Time
	closed bool
}

func (c *channel) close() {
	c.ch.Close()
	c.closed = true
}