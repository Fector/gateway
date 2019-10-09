package router

import (
	"github.com/DeathHand/gateway/callback"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
)

type Router struct {
	gateways  *[]model.Gateway
	ingress   *chan model.Message
	egress    *chan model.Message
	memory    memory.Memory
	connector *Connector
	callback  callback.Callback
	error     *chan error
}

func NewRouter(
	gateways *[]model.Gateway,
	ingress *chan model.Message,
	egress *chan model.Message,
	memory memory.Memory,
	error *chan error,
) *Router {
	return &Router{
		gateways:  gateways,
		ingress:   ingress,
		egress:    egress,
		memory:    memory,
		connector: NewConnector(gateways, egress, error),
		callback:  callback.NewHttpCallback(error),
		error:     error,
	}
}

func (r *Router) Ingress() {
}

func (r *Router) Egress() {
	message := <-*r.egress
	r.callback.Add(&message)
}

func (r *Router) Run() {
	go r.Egress()
	go r.Ingress()
	go r.callback.Run()
	go r.connector.Run()
}
