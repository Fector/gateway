package storage

import (
	"github.com/DeathHand/gateway/model"
)

type Message struct {
	data map[string]*model.Message
}

func (m *Message) AddMessage(c *model.Message) {

}

func (m *Message) DeleteMessage(uuid string) {

}

func (m *Message) CollectExpired() {

}
