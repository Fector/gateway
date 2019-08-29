package router

import (
	"fmt"
	"github.com/DeathHand/gateway/callback"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"log"
)

type Router struct {
	ingress  map[string]*chan model.Message
	egress   map[string]*chan model.Message
	memory   memory.Memory
	callback callback.Callback
}

func NewRouter(
	ingress map[string]*chan model.Message,
	egress map[string]*chan model.Message,
	memory memory.Memory,
) *Router {
	return &Router{
		ingress: ingress,
		egress:  egress,
		memory:  memory,
	}
}

func (r *Router) MoIngress() {
	m := <-*r.ingress["mo"]
	c, exists := r.egress[m.Gateway]
	if !exists {
		log.Println(fmt.Sprintf("%s gateway does not exist or not connected", m.Gateway))
	}
	err := r.memory.Put(&m)
	if err != nil {
		log.Println(err)
	}
	*c <- m
}

func (r *Router) MtIngress() {
	m := <-*r.ingress["mt"]
	c, exists := r.egress[m.Gateway]
	if !exists {
		log.Println(fmt.Sprintf("%s gateway does not exist or not connected", m.Gateway))
	}
	err := r.memory.Put(&m)
	if err != nil {
		log.Println(err)
	}
	*c <- m
}

func (r *Router) Notify() {
	m := <-*r.memory.Notify()
	r.callback.Send(&m)
	log.Println(m.Uuid)
}

func (r *Router) Run() {
	go r.Notify()
	go r.MtIngress()
	go r.MoIngress()
	go r.memory.Observe()
}
