package go_rabbitmq

import (
	"github.com/streadway/amqp"
	"time"
)

type connection struct {
	conn *amqp.Connection
	createdAt time.Time
	closed bool
	channelList []*channel
}

func (c *connection) close() {
	for channelIndex := range c.channelList {
		c.channelList[channelIndex].close()
	}
	c.conn.Close()
	c.channelList = nil
	c.closed = true
}