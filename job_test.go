package go_rabbitmq

import (
	"testing"
	"time"
)

var defaultConfig = &Config{
	Username: "golang",
	Password: "123456",
	VirtualHost: "go_vhost",
	ConnCap: 2,
	ChCap: 5,
}

type OrderJob struct {
	
}

func (o OrderJob) RoutingKey() string {
	return o.Queue().Binding.BindingKey
}

func (o OrderJob) Exchange() *Exchange {
	return nil
}

func (o OrderJob) Queue() *Queue {
	return &Queue{
		Name: "Order",
		Binding: QueueBind{
			BindingKey: "shop.order.submit",
		},
	}
}

func (o OrderJob) Handle(payload Payload) {
	panic("implement me")
}

func (o OrderJob) Config() *Config {
	return defaultConfig
}

func TestDispatch(t *testing.T) {
	err := Dispatch(&OrderJob{}, Payload{
		"order_id": "100001",
		"price": 10.35,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("push done")
	}
}

func TestLotsDispatch(t *testing.T) {
	for i := 0; i < 100; i++ {
		<-time.After(time.Second)
		TestDispatch(t)
	}
}
