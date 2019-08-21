package engine

import "github.com/DeathHand/gateway/model"

type Engine interface {
	Put(message *model.Message) error
	Get(uuid string) (*model.Message, error)
	Delete(uuid string) error
	Run() error
}
