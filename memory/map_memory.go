package memory

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"github.com/google/uuid"
	"log"
	"os"
	"sync"
	"time"
)

type MapMemory struct {
	Memory
	data   map[string]*model.Message
	path   string
	mutex  sync.Mutex
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

	return &MapMemory{path: path}, nil
}

func (m *MapMemory) Put(message *model.Message) (string, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	id, err := uuid.NewUUID()

	if err != nil {
		return "", err
	}

	m.data[id.String()] = message

	return id.String(), nil
}

func (m *MapMemory) Get(id string) (*model.Message, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	if message, exists := m.data[id]; exists {
		return message, nil
	}

	return nil, errors.New(fmt.Sprintf("%s not found ", id))
}

func (m *MapMemory) Delete(id string) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	if _, exists := m.data[id]; exists {
		delete(m.data, id)
	}

	return errors.New(fmt.Sprintf("%s not found ", id))
}

func (m *MapMemory) Collect() {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	for id, message := range m.data {
		if time.Now().UnixNano() >= message.Timestamp+message.Ttl {
			log.Println(fmt.Sprintf("%s expired", id))
			delete(m.data, id)
		}
	}
}

func (m *MapMemory) Dump() {
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

func (m *MapMemory) Observe() error {
	dump := time.NewTicker(10 * time.Second)
	collect := time.NewTicker(30 * time.Second)

	select {
	case <-dump.C:
		go m.Dump()
	case <-collect.C:
		go m.Collect()
	}

	return nil
}
