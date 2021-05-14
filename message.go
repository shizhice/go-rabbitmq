package go_rabbitmq

import (
	"encoding/json"
	"log"
	"time"
)

type messageState int8

const (
	messageSend messageState = iota
	messageReceive
	messageNotMatchQueue
	messageDead
)

type messageId string

type message struct {
	Id messageId `json:"id"`
	Address addressHash `json:"address"`
	ExchangeName string `json:"exchange_name"`
	Queue string `json:"queue"`
	Retry int `json:"retry"`
	Payload Payload `json:"payload"`
	State messageState `json:"state"`
	DeadAt time.Time `json:"dead_at"`
	CreateAt time.Time `json:"create_at"`
}

func (m message) Marshal() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Printf("payload marshal error: %v", err)
	}

	return bytes
}


type Payload map[string]interface{}

