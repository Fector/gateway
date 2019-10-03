package connection

import (
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"log"
)

type Connector struct {
	gateways    []model.Gateway
	connections map[string]*Connection
	Error       *chan error
}

func NewConnector(gateways []model.Gateway, error *chan error) *Connector {
	return &Connector{
		gateways: gateways,
		Error:    error,
	}
}

func (c *Connector) AddConnection(gateway *model.Gateway) error {
	stop := make(chan int, 1)
	inbox := make(chan model.Message, 65535)
	outbox := make(chan model.Message, 65535)
	connection, err := NewConnection(gateway, &inbox, &outbox, &stop, c.Error)
	if err != nil {
		return err
	}
	c.connections[gateway.Name] = connection
	go c.connections[gateway.Name].Run()
	return nil
}

func (c *Connector) StopConnection(gateway string) error {
	if connection, exists := c.connections[gateway]; exists {
		connection.Stop()
		return nil
	}
	return errors.New(fmt.Sprintf("No %s connection found ", gateway))
}

func (c *Connector) RemoveConnection(gateway string) error {
	if connection, exists := c.connections[gateway]; exists {
		connection.Stop()
		return nil
	}
	return errors.New(fmt.Sprintf("No %s connection found ", gateway))
}

func (c *Connector) SendMessage(gateway string, message *model.Message) error {
	if connection, exists := c.connections[gateway]; exists {
		connection.SendMessage(message)
		return nil
	}
	return errors.New(fmt.Sprintf("No %s connection found ", gateway))
}

func (c *Connector) Run() {
	for _, gate := range c.gateways {
		err := c.AddConnection(&gate)
		if err != nil {
			log.Fatal(err)
		}
	}
}
