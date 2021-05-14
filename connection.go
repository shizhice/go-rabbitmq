package go_rabbitmq

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type connection struct {
	conn *amqp.Connection
	createdAt time.Time
	closed bool
	channelList []*channel
	quit chan bool
	closeChan chan *amqp.Error
}

func (c *connection) close() {
	if c.closed {
		return
	}
	if ! c.conn.IsClosed() {
		for channelIndex := range c.channelList {
			c.channelList[channelIndex].close()
		}
		c.conn.Close()
	}

	defer func() {
		c.channelList = nil
		c.closed = true
		c.quit <- true
	}()
}

func (c *connection) listen() {
	c.quit = make(chan bool)
	c.closeChan = c.conn.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-c.quit:
			return
		case err := <-c.closeChan:
			if !c.closed {
				log.WithFields(log.Fields{
					"error": err,
				}).Info("connection 关闭")
				c.close()
			}
		}
	}
}