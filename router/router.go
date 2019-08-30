package router

import (
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/gateway/service"
)

type Router struct {
	ingress  *chan model.Message
	egress   *chan model.Message
	callback *chan model.Message
	memory   memory.Memory
	errChan  *chan error
}

func NewRouter(s *service.Service) *Router {
	return &Router{
		ingress: s.Ingress,
		egress:  s.Egress,
		memory:  s.Memory,
		errChan: s.ErrChan,
	}
}

func (r *Router) Callback() {
	m := <-*r.memory.Notify()
	*r.callback <- m
}

func (r *Router) Run() {
	go r.Callback()
}
