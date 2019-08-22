package memory

import "github.com/DeathHand/gateway/model"

type Memory interface {
	Put(message *model.Message) (string, error)
	Get(id string) (*model.Message, error)
	Delete(id string) error
	Observe()
}
