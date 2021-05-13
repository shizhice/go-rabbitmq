package go_rabbitmq

import (
	"github.com/streadway/amqp"
)

type exchangeType int8

const (
	DirectExchange exchangeType = iota
	FanoutExchange
	TopicExchange
	HeaderExchange
	maxExchangeType
)

func (e *exchangeType) Name() (name string) {
	switch *e {
	case 0: name = "direct"
	case 1: name = "fanout"
	case 2: name = "topic"
	case 3: name = "headers"
	}

	return
}

type Exchange struct {
	// 交换机名字
	Name string

	// 交换机类型
	Kind exchangeType

	// 扩展参数
	arguments  amqp.Table
}

func (e *Exchange) declare(ch *amqp.Channel) error {
	//   ------------------------------------------------------------------------------------------------------------
	//   |      参数     |    类型    |                              参数说明                                          |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  name        |    string  |  交换机名字                                                                    |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  kind        |    string  |  交换机类型                                                                    |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  durable     |    bool    |  是否持久化，建议为true                                                          |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  autoDelete  |    bool    |  当exchange绑定的queue全都解绑的时候exchange会自动删除， 建议为 false                |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  internal    |    bool    |  内部交换机，只能通过管理后台发送消息, 建议为 false                                  |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  noWait      |    bool    |  声明了一个交换器之后, 是否需要等待服务器的返回, 通常设置为false. 当设置为true时，声明一   |
	//   |              |            |  个交换器后，实际RabbitMQ还未完成此交换器的创建, 那么此时向此交换器投递消息，就会发生异常。  |
	//   ------------------------------------------------------------------------------------------------------------
	//   |  args        | amqp.Table |  扩展参数                                                                      |
	//   -----------------------------------------------------------------------------------------------------------

	return ch.ExchangeDeclare(e.Name, e.Kind.Name(), true, false, false, false, e.arguments)
}

var commonDirectExchange = &Exchange{
	Name: commonDirectExchangeName,
	Kind: DirectExchange,
}

var commonFanoutExchange = &Exchange{
	Name: commonFanoutExchangeName,
	Kind: FanoutExchange,
}

var commonTopicExchange = &Exchange{
	Name: commonTopicExchangeName,
	Kind: TopicExchange,
}

var commonHeaderExchange = &Exchange{
	Name: commonHeaderExchangeName,
	Kind: HeaderExchange,
}

var deadLetterExchange = &Exchange{
	Name: deadLetterExchangeName,
	Kind: DirectExchange,
}
