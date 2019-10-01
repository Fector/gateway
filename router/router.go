package router

import (
	"github.com/DeathHand/gateway/callback"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/gateway/service"
)

type Router struct {
	ingress  *chan model.Message
	egress   *chan model.Message
	memory   memory.Memory
	callback callback.Callback
	error    *chan error
}

func NewRouter(s *service.Service) *Router {
	return &Router{
		ingress:  s.Ingress,
		egress:   s.Egress,
		memory:   s.Memory,
		callback: callback.NewHttpCallback(s.Error),
		error:    s.Error,
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
}
