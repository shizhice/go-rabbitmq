package go_rabbitmq

import (
	"crypto/md5"
	"fmt"
)

type addressHash string
type address string

func (a address) md5() addressHash {
	return addressHash(fmt.Sprintf("%x", md5.Sum([]byte(a))))
}

type Config struct {
	addr address
	Username string
	Password string
	Host string
	Port int
	VirtualHost string
	ConnCap int
	ChannelCapOfPreCoon int
}

func (c *Config) init() {
	if c.Username == "" {
		c.Username = "guest"
	}
	if c.Password == "" {
		c.Password = "guest"
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == 0 {
		c.Port = 5672
	}
	if c.VirtualHost == "" {
		c.VirtualHost = "/"
	}
	if c.ConnCap == 0 {
		c.ConnCap = 1
	}
	if c.ChannelCapOfPreCoon == 0 {
		c.ChannelCapOfPreCoon = 1
	}

	c.addr = address(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", c.Username, c.Password, c.Host, c.Port, c.VirtualHost))
}

func (c *Config) getAddr() address {
	if c.addr == "" {
		c.init()
	}

	return c.addr
}