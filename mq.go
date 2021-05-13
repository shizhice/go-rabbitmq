package go_rabbitmq

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

var mqCollect = rabbitmqCollect{mqMap: make(map[address]*MQ), queueMap: make(map[string]bool), exchangeMap: make(map[string]bool)}

type rabbitmqCollect struct {
	sync.Mutex
	mqMap map[address]*MQ
	queueMu sync.Mutex
	queueMap map[string]bool
	exchangeMu sync.Mutex
	exchangeMap map[string]bool
}

// Open a RabbitMQ connection pool and return MQ pointer
func Open(conf Config) (mq *MQ, err error) {
	mqCollect.Lock()
	defer mqCollect.Unlock()
	conf.init()

	if _, ok := mqCollect.mqMap[conf.addr]; ! ok {
		mq, err := newMQ(conf)
		if err != nil {
			return nil, err
		}

		mqCollect.mqMap[conf.addr] = mq
	}

	return mqCollect.mqMap[conf.addr], nil
}

func newMQ(conf Config) (*MQ, error) {
	pool, err := CreatePool(conf.addr, conf.ConnCap, conf.ChCap)
	if err != nil {
		return nil, err
	}
	mq := &MQ{conf: conf, pool: pool}

	return mq, nil
}

type MQ struct {
	conf Config
	pool *pool
}

func (mq *MQ) Close() {
	mq.pool.close()
}

// SetExchange set connection`s exchange
func (mq *MQ) DeclareExchange(exchange ...*Exchange) error {
	mqCollect.exchangeMu.Lock()
	defer mqCollect.exchangeMu.Unlock()
	activeCh, err := mq.pool.get()
	defer mq.pool.release(activeCh)
	if err != nil {
		return err
	}
	if len(exchange) > 0 {
		for i := range exchange {
			if _, ok := mqCollect.exchangeMap[exchange[i].Name]; !ok {
				if err := exchange[i].declare(activeCh.ch); err != nil {
					return err
				}
				mqCollect.exchangeMap[exchange[i].Name] = true
			}
		}
	}

	return nil
}

func (mq *MQ) DeclareQueue(queue ...*Queue) error {
	mqCollect.queueMu.Lock()
	defer mqCollect.queueMu.Unlock()
	activeCh, err := mq.pool.get()
	defer mq.pool.release(activeCh)
	if err != nil {
		return err
	}
	if len(queue) > 0 {
		for i := range queue {
			if _, ok := mqCollect.queueMap[queue[i].Name]; !ok {
				if err := queue[i].declare(activeCh.ch); err != nil {
					return err
				}
				mqCollect.queueMap[queue[i].Name] = true
			}
		}
	}
	return nil
}

// get the MQ`s pool length
func (mq *MQ) PoolLen() (chNum int, connNum int) {
	return mq.pool.len()
}