package memory

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"os"
	"sync"
	"time"
)

type MapMemory struct {
	Memory
	data   map[string]model.Message
	path   string
	mutex  sync.Mutex
	dump   *time.Ticker
	expire *time.Ticker
	notify *chan model.Message
	errors *chan error
}

func NewMemoryEngine(path string) (*MapMemory, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("Message engine path is not a directory ")
	}
	return &MapMemory{
		path:   path,
		dump:   time.NewTicker(10 * time.Second),
		expire: time.NewTicker(30 * time.Second),
	}, nil
}

func (m *MapMemory) Put(message model.Message) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.data[message.Uuid] = message
	return nil
}

func (m *MapMemory) Get(uuid string) (*model.Message, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if message, exists := m.data[uuid]; exists {
		return &message, nil
	}
	return nil, errors.New(fmt.Sprintf("%s not found ", uuid))
}

func (m *MapMemory) Delete(uuid string) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if _, exists := m.data[uuid]; exists {
		delete(m.data, uuid)
	}
	return errors.New(fmt.Sprintf("%s not found ", uuid))
}

func (m *MapMemory) collect() {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	for id, message := range m.data {
		if time.Now().UnixNano() >= message.Timestamp+message.Ttl {
			delete(m.data, id)
			*m.notify <- message
		}
	}
}

func (m *MapMemory) save() {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(m.data)
	if err != nil {
		*m.errors <- err
	}
}

func (m *MapMemory) Restore(b *bytes.Buffer) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	dec := gob.NewDecoder(b)
	err := dec.Decode(m.data)
	if err != nil {
		return err
	}
	return nil
}

func (m *MapMemory) Observe() {
	for {
		select {
		case <-m.dump.C:
			m.save()
		case <-m.expire.C:
			m.collect()
		}
	}
}

func (m *MapMemory) Notify() *chan model.Message {
	return m.notify
}
