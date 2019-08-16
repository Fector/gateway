package storage

import (
	"github.com/DeathHand/gateway/model"
)

type Message struct {
	data map[string]*model.Message
}

func (m *Message) Get(uuid string) {

}

func (m *Message) Add(c *model.Message) {

}

func (m *Message) Delete(uuid string) {

}

func (m *Message) Collect() {

}
