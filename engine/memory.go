package engine

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"log"
	"os"
	"sync"
	"time"
)

type MemoryEngine struct {
	Engine
	data   map[string]*model.Message
	path   string
	mutex  sync.Mutex
	router *Router
	errors *chan error
}

func NewMemoryEngine(path string) (*MemoryEngine, error) {
	info, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("Message engine path is not a directory ")
	}

	return &MemoryEngine{path: path}, nil
}

func (m *MemoryEngine) Get(uuid string) (*model.Message, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	return &model.Message{}, nil
}

func (m *MemoryEngine) Put(message *model.Message) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	return nil
}

func (m *MemoryEngine) Delete(uuid string) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	return nil
}

func (m *MemoryEngine) collect() {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	for uuid, message := range m.data {
		if time.Now().UnixNano() >= message.Timestamp+message.Ttl {
			log.Println(fmt.Sprintf("%s expired", uuid))
			delete(m.data, uuid)
		}
	}
}

func (m *MemoryEngine) dump() {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	err := enc.Encode(m.data)

	if err != nil {
		*m.errors <- err
	}
}

func (m *MemoryEngine) Restore(b *bytes.Buffer) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	d := gob.NewDecoder(b)

	err := d.Decode(m.data)

	if err != nil {
		return err
	}

	return nil
}

func (m *MemoryEngine) Run() error {
	dumpTicker := time.NewTicker(10 * time.Second)
	collectTicker := time.NewTicker(30 * time.Second)

	select {
	case <-dumpTicker.C:
		go m.dump()
	case <-collectTicker.C:
		go m.collect()
	}

	return nil
}
