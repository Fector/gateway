package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/DeathHand/gateway/model"
	"os"
	"sync"
)

type Message struct {
	data  map[string]*model.Message
	path  string
	mutex sync.Mutex
}

func NewMessageStorage(path string) (*Message, error) {
	info, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("Message storage path is not a directory ")
	}

	return &Message{path: path}, nil
}

func (m *Message) Get(uuid string) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
}

func (m *Message) Put(c *model.Message) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
}

func (m *Message) Delete(uuid string) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
}

func (m *Message) Collect() {
	defer m.mutex.Unlock()
	m.mutex.Lock()
}

func (m *Message) Backup() {
	defer m.mutex.Unlock()
	m.mutex.Lock()
}

func (m *Message) dump() (*bytes.Buffer, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	var b bytes.Buffer

	e := gob.NewEncoder(&b)

	err := e.Encode(m.data)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (m *Message) restore(b *bytes.Buffer) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	d := gob.NewDecoder(b)

	err := d.Decode(m.data)

	if err != nil {
		return err
	}

	return nil
}
