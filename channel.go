package go_rabbitmq

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type channel struct {
	ch *amqp.Channel
	createdAt time.Time
	closed bool
	closeChan chan *amqp.Error
	notifyConfirmChan chan amqp.Confirmation
	notifyReturnChan  chan amqp.Return
	publishChan chan message
	quit chan bool
}

func (c *channel) close() {
	if c.closed {
		return
	}
	defer c.ch.Close()
	c.closed = true
	c.quit <- true
}

func (c *channel) listen() {
	var deliveryTag uint64 = 1
	var deliveryMap = make(map[uint64]messageId)
	c.quit = make(chan bool)
	c.publishChan = make(chan message)
	c.notifyConfirmChan = c.ch.NotifyPublish(make(chan amqp.Confirmation))
	c.notifyReturnChan = c.ch.NotifyReturn(make(chan amqp.Return))
	c.closeChan = c.ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-c.quit:
			return
		case msg := <-c.publishChan:
			deliveryMap[deliveryTag] = msg.Id
			deliveryTag++
		case confirmed := <-c.notifyConfirmChan:
			if confirmed.Ack {
				if msgId, ok := deliveryMap[confirmed.DeliveryTag]; ok {
					log.WithFields(log.Fields{
						"MsgId": msgId,
						"DeliveryTag": confirmed.DeliveryTag,
					}).Info("消息确认送达")
					delete(deliveryMap, confirmed.DeliveryTag)
				}
			}
		case ret := <-c.notifyReturnChan:
			c.notifyReturn(ret)
		case err := <-c.closeChan:
			if !c.closed {
				c.closed = true
				log.WithFields(log.Fields{
					"error": err,
				}).Info("channel 关闭")
			}
		}
	}
}

func (c *channel) notifyReturn(ret amqp.Return) {
	if string(ret.Body) != "" {
		var msg = &message{}
		_ = json.Unmarshal(ret.Body, msg)
		log.WithFields(log.Fields{
			"MessageID": msg.Id,
		}).Warn("消息没有正确入列")
	}
}

func (c *channel) publish(exchangeName , routingKey string, msg message, ttl TTL) error {
	if confirmErr := c.ch.Confirm(false); confirmErr != nil {
		return confirmErr
	}

	defer func() {
		c.publishChan <- msg
		log.WithFields(log.Fields{
			"ExchangeName": exchangeName,
			"RoutingKey": routingKey,
			"MessageID": msg.Id,
		}).Info("发送消息")
	}()

	return c.ch.Publish(exchangeName, routingKey, true, false, c.makeMessage(msg, ttl))
}

func (c *channel) makeMessage(m message, ttl TTL) amqp.Publishing {
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        m.Marshal(),
		DeliveryMode: 2,
	}

	// 设置message过期时间
	if ttl >= MixTTL {
		msg.Expiration =  ttl.String()
	}

	return msg
}