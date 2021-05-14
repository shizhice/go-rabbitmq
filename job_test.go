package go_rabbitmq

import (
	"testing"
	"time"
)

var defaultConfig = &Config{
	Username: "golang",
	Password: "123456",
	VirtualHost: "go_vhost",
	//ConnCap: 2,
	//ChCap: 5,
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
	}
}

func TestDispatchNotConfirm(t *testing.T) {
	for i := 0; i < 100; i++ {
		TestDispatch(t)
	}
	time.Sleep(time.Second * 5)
}

func TestLotsDispatch(t *testing.T) {
	for i := 0; i < 100; i++ {
		<-time.After(time.Second)
		TestDispatch(t)
	}
}

func TestConcurrentDispatch(t *testing.T) {
	for n := 0; n < 10; n++ {
		go func() {
			for i := 0; i < 10; i++ {
				TestDispatch(t)
			}
		}()
	}
	time.Sleep(time.Second * 2)
}

type NotifyReturnJob struct {

}

func (n NotifyReturnJob) Config() *Config {
	return defaultConfig
}

func (n NotifyReturnJob) Queue() *Queue {
	return &Queue{
		Name: "NotifyReturnQueue",
		Binding: QueueBind{
			BindingKey: "queue.notifyReturn",
		},
	}
}

func (n NotifyReturnJob) Exchange() *Exchange {
	return nil
}

func (n NotifyReturnJob) Handle(payload Payload) {
	panic("implement me")
}

func (n NotifyReturnJob) RoutingKey() string {
	return "queue.notifyReturn.test"
}

func TestNotifyReturn(t *testing.T) {
	err := Dispatch(&NotifyReturnJob{}, Payload{
		"message": "test",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestLostNotifyReturn(t *testing.T) {
	for i := 0; i < 100; i++ {
		TestNotifyReturn(t)
	}
}

func TestConcurrentNotifyReturn(t *testing.T) {
	for n := 0; n < 10; n++ {
		go func() {
			for i := 0; i < 10; i++ {
				TestNotifyReturn(t)
			}
		}()
	}
	time.Sleep(time.Second * 5)
}