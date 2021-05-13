package go_rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

// 连接池
type pool struct {
	mu sync.Mutex    // 互斥锁，保证编程安全
	addr address
	connCap int // connection 连接数
	chCap int // 每个connection创建的channel，totalChannel = connCap * channelCap
	closed  bool
	connList []*connection
	readyChannel chan *channel
}

// CreatePool crate a rabbitmq pool
func CreatePool(addr address, connCap, chCap int) (*pool, error) {
	if connCap == 0 {
		connCap = 1
	}
	if chCap == 0 {
		chCap = 1
	}
	return (&pool{
		addr: addr,
		connCap: connCap,
		chCap: chCap,
		readyChannel: make(chan *channel, connCap * chCap),
	}).create()
}

// Close the rabbitmq pool
func (p *pool) close() {
	if p.closed {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = true
	close(p.readyChannel)

	for connIndex := range p.connList {
		p.connList[connIndex].close()
	}
	p.connList = nil
}

// Get a rabbitmq channel
func (p *pool) get() (*channel, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}

	activeChannel, ok := <-p.readyChannel
	if !ok {
		return nil, ErrPoolClosed
	}

	return activeChannel, nil
}

// Release the rabbitmq channel
func (p *pool) release(ch *channel) error {
	if p.closed {
		return errors.New("pool closed")
	}

	select {
	case p.readyChannel <- ch:
		return nil
	default:
		return errors.New("ConnPool is full")
	}
}

func (p *pool) create() (*pool, error) {
	p.mu.Lock()
	err := p.createConn()
	p.mu.Unlock()
	// 创建连接失败将关闭连接池
	if err != nil {
		p.close()
	}
	return p, err
}

func (p *pool) createConn() (err error) {
	for i := 0; i < p.connCap; i++ {
		conn, err := amqp.Dial(string(p.addr))
		if err != nil {
			return err
		}
		curConn := &connection{
			conn: conn,
			createdAt: time.Now(),
		}
		p.connList = append(p.connList, curConn)
		err = p.createChannel(curConn)
		if err != nil {
			return err
		}
	}

	return nil
}

// createChannel 创建connection`s channel
func (p *pool) createChannel(conn *connection) (err error) {
	for i := 0; i < p.chCap; i++ {
		ch, err := conn.conn.Channel()
		if err != nil {
			return err
		}
		curChannel := &channel{
			ch: ch,
			createdAt: time.Now(),
		}
		conn.channelList = append(conn.channelList, curChannel)
		p.readyChannel <- curChannel
	}

	return nil
}

// createChannel 创建connection`s channel
func (p *pool) len() (chNum int, connNum int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for connIndex := range p.connList {
		fmt.Println(p.connList[connIndex].channelList)
		chNum += len(p.connList[connIndex].channelList)
	}

	return chNum, len(p.connList)
}