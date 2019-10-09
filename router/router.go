package router

import (
	"github.com/DeathHand/gateway/callback"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"log"
	"os"
)

type Router struct {
	gateways  *[]model.Gateway
	ingress   *chan model.Message
	egress    *chan model.Message
	memory    memory.Memory
	connector *Connector
	callback  callback.Callback
	logger    *log.Logger
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
		logger: log.New(
			os.Stdout,
			"Router: ",
			log.Ldate|log.Ltime|log.Lmicroseconds,
		),
		error: error,
	}
}

func (r *Router) Ingress() {
	for {
		message := <-*r.ingress
		conn, err := r.connector.GetConnection(message.Gateway)
		if err != nil {
			r.logger.Printf("Message (uid: %s) error: destination gateway (%s) not found.", message.Uuid, message.Gateway)
			*r.egress <- message
			continue
		}
		*conn.Inbox <- message
	}
}

func (r *Router) Egress() {
	for {
		message := <-*r.egress
		r.callback.Add(&message)
	}
}

func (r *Router) Run() {
	go r.Egress()
	go r.Ingress()
	go r.callback.Run()
	go r.connector.Run()
}
