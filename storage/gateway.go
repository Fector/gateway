package storage

import (
	"github.com/DeathHand/gateway/model"
)

type Gateway struct {
	data map[string]*model.Gateway
}

func (g *Gateway) Get(uuid string) {

}

func (g *Gateway) Add(c *model.Gateway) string {
	return ""
}

func (g *Gateway) Delete(uuid string) {

}
