package memory

import "github.com/DeathHand/gateway/model"

/**
Memory is a memory storage interface.

Implements Put() function to store model.Message data.
Implements Get() function to get model.Message data.
Implements Delete() function to delete model.Message data.
Implements Observe() function to search and destroy expired messages,
and notify Router about that event.

When implementing new memory, you must implement this methods
and observe memory data structures for messages,
which time.Now().UnixNano() >= message.Timestamp+message.Ttl
*/
type Memory interface {
	Put(message *model.Message) error
	Get(uuid string) (*model.Message, error)
	Delete(uuid string) error
	Notify() *chan model.Message
	Run()
}
