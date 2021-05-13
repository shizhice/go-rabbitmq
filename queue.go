package go_rabbitmq

import (
	"github.com/streadway/amqp"
)

// Queue 队列
type Queue struct {
	Name string
	Binding QueueBind
	Arguments amqp.Table
}

func (q *Queue) isDelayQueue() bool {
	if q.Arguments == nil {
		return false
	}

	for key := range q.Arguments {
		if key == "x-dead-letter-exchange" || key == "x-dead-letter-queue" || key == "x-message-ttl" {
			return true
		}
	}

	return false
}

func (q *Queue) declare(ch *amqp.Channel) error {
	//   ------------------------------------------------------------------------------------------------------------
	//   |      参数     |    类型    |                              参数说明                                          |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  name        |    string  |  路由名字                                                                      |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  durable     |    bool    |  是否持久化，建议为true                                                          |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  autoDelete  |    bool    |  当queue无消息时会自动删除， 建议为 false                                          |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  exclusive   |    bool    |  排他队列，只对声明该队列的用户可见，其它用户无法访问。                                 |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  noWait      |    bool    |  声明了一个Queue之后, 是否需要等待服务器的返回, 通常设置为false. 当设置为true时，声明一   |
	//   |              |            |  个Queue后，实际RabbitMQ还未完成此Queue的创建, 那么此时向此Queue投递消息，就会发生失败。 |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  args        | amqp.Table |  扩展参数                                                                      |
	//   |              |            |  常用参数:                                                                     |
	//   |              |            |     ① x-message-ttl 该队列内消息的存活时间（毫秒）。                               |
	//   |              |            |     ② x-dead-letter-exchange 队列溢出、消息被拒或过期 重定向的路由。                |
	//   |              |            |     ③ x-dead-letter-exchange 队列溢出、消息被拒或过期 重定向的队列。                |
	//   |              |            |     ④ x-dead-letter-routing-key 当 队列溢出、消息被拒或过期，重定向时替换Message的   |
	//   |              |            |       RoutingKey                                                              |
	//   ------------------------------------------------------------------------------------------------------------
	queue, err := ch.QueueDeclare(q.Name, true, false, false, false, q.Arguments)
	if err == nil {
		err = ch.QueueBind(
			queue.Name,
			q.Binding.BindingKey,
			q.Binding.Exchange.Name,
			false,
			q.Binding.Args)
	}

	return err

}

// QueueBind binding key
type QueueBind struct {
	Exchange *Exchange
	BindingKey string
	Args amqp.Table
}
// todo: 不设置 BindingKey message会路由到队列吗？

var deadLetterQueue = &Queue{
	Name: deadLetterQueueName,
	Binding: QueueBind{
		Exchange: deadLetterExchange,
	},
	Arguments: nil,
}