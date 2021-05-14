package go_rabbitmq

import (
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	t.Log("实例化MQ.")
	mq, err := Open(Config{
		Username: "golang",
		Password: "123456",
		VirtualHost: "go_vhost",
		ConnCap: 2,
		ChannelCapOfPreCoon: 5,
	})

	if err != nil {
		t.Errorf("open mq error: %v", err)
	} else {
		chNum, connNum := mq.PoolLen()
		t.Logf("实例化MQ完成. channel num: %d ; connection num: %d", chNum, connNum)
		time.Sleep(time.Second * 60)
	}
}
