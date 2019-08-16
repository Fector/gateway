package storage

import (
	"github.com/DeathHand/gateway/model"
)

type Gateway struct {
	data map[string]*model.Gateway
}

func (g *Gateway) GetGateway(uuid string) {

}

func (g *Gateway) AddGateway(c *model.Gateway) string {
	return ""
}

func (g *Gateway) DeleteGateway(uuid string) {

}
