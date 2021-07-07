package go_rabbitmq

import (
	"github.com/streadway/amqp"
	"testing"
	"time"
)

var defaultConfig = &Config{
	Username: "golang",
	Password: "123456",
	VirtualHost: "go_vhost",
	ConnCap: 1,
	ChannelCapOfPreCoon: 1,
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
	//time.Sleep(time.Second * 10)
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

type CapacityJob struct {

}

func (c CapacityJob) Config() *Config {
	return defaultConfig
}

func (c CapacityJob) Queue() *Queue {
	return &Queue{
		Name: "CapacityQueue",
		Binding: QueueBind{
			BindingKey: "capacityQueue",
		},
		Arguments: amqp.Table{
			"x-max-length": 10,
			"x-overflow": "drop-head",
			//* drop-head（删除queue头部的消息）、
			//* reject-publish-dlx（拒绝发送消息到死信交换器）
			//* 类型为quorum 的queue只支持drop-head;
		},
	}
}

func (c CapacityJob) Exchange() *Exchange {
	return nil
}

func (c CapacityJob) Handle(payload Payload) {
	panic("implement me")
}

func (c CapacityJob) RoutingKey() string {
	return c.Queue().Binding.BindingKey
}

func TestDispatchWithCapacity(t *testing.T) {
	err := Dispatch(CapacityJob{}, Payload{
		"id": 10001 + num,
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Log("push done")
	}

}

var num = 0
func TestLotsDispatchWithCapacity(t *testing.T) {
	for i := 0; i < 20; i++ {
		<-time.After(time.Millisecond * 1000)
		TestDispatchWithCapacity(t)
		num++
	}

	time.Sleep(time.Second * 2)
}