package service

import (
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/conn"
	"github.com/DeathHand/gateway/model"
)

type Connector struct {
	connections map[string]*conn.Conn
	Error       *chan error
}

func (c *Connector) AddConnection(gateway *model.Gateway) error {
	stop := make(chan int, 1)
	inbox := make(chan model.Message, 65535)
	outbox := make(chan model.Message, 65535)
	connection, err := conn.NewConnection(gateway, &inbox, &outbox, &stop, c.Error)
	if err != nil {
		return err
	}
	c.connections[gateway.Name] = connection
	return nil
}

func (c *Connector) RunConnection(gateway string) error {
	if connection, exists := c.connections[gateway]; exists {
		go connection.Run()
		return nil
	}
	return errors.New(fmt.Sprintf("No %s connection found ", gateway))
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

func (c *Connector) Collect() {

}
