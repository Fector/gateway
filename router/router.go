package router

import (
	"fmt"
	"github.com/DeathHand/gateway/callback"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"log"
)

type Router struct {
	Ingress  map[string]*chan model.Message
	Egress   map[string]*chan model.Message
	Memory   memory.Memory
	CallBack callback.Callback
}

// AddEgress is a new connection inbox setter
func (r *Router) AddEgress(id string, e *chan model.Message) {
	r.Egress[id] = e
}

func (r Router) MoIngress() {
	m := <-*r.Ingress["mo"]
	c, exists := r.Egress[m.Gateway]
	if !exists {
		log.Println(fmt.Sprintf("%s gateway does not exist or not connected", m.Gateway))
	}
	err := r.Memory.Put(&m)
	if err != nil {
		log.Println(err)
	}
	*c <- m
}

func (r Router) MtIngress() {
	m := <-*r.Ingress["mt"]
	c, exists := r.Egress[m.Gateway]
	if !exists {
		log.Println(fmt.Sprintf("%s gateway does not exist or not connected", m.Gateway))
	}
	err := r.Memory.Put(&m)
	if err != nil {
		log.Println(err)
	}
	*c <- m
}

func (r *Router) Notify() {
	m := <-*r.Memory.Notify()
	r.CallBack.Send(&m)
	log.Println(m.Uuid)
}

func (r *Router) Run() {
	go r.Notify()
	go r.MtIngress()
	go r.MoIngress()
}
