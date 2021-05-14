package go_rabbitmq

import (
	"testing"
	"time"
)


type AutoCancelOrderJob struct {

}

func (a AutoCancelOrderJob) TTL() TTL {
	return TTLSecond * 10
}

func (a AutoCancelOrderJob) DlxExchange() *Exchange {
	return nil
}

func (a AutoCancelOrderJob) DlxQueue() *Queue {
	return nil
}

func (a AutoCancelOrderJob) RoutingKey() string {
	return a.Queue().Binding.BindingKey
}

func (a AutoCancelOrderJob) Exchange() *Exchange {
	return nil
}

func (a AutoCancelOrderJob) Queue() *Queue {
	return &Queue{
		Name: "autoCancelOrder",
		Binding: QueueBind{
			BindingKey: "shop.order.cancel.auto",
		},
	}
}

func (a AutoCancelOrderJob) Handle(payload Payload) {
	panic("implement me")
}

func (a AutoCancelOrderJob) Config() *Config {
	return defaultConfig
}

func TestDispatchDelay(t *testing.T) {
	err := DispatchDelay(&AutoCancelOrderJob{}, Payload{
		"order_id": "100001",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestLotsDispatchDelay(t *testing.T) {
	for i := 0; i < 100; i++ {
		<-time.After(time.Second)
		TestDispatchDelay(t)
	}
}