package router

import (
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/connection"
	"github.com/DeathHand/gateway/model"
	"log"
	"time"
)

// Connector represents connection pool
// This type has functions to run all configured SMPP Gateway pool
type Connector struct {
	gateways *[]model.Gateway
	egress   *chan model.Message
	pool     map[string]*connection.Connection
	timer    *time.Timer
	Error    *chan error
}

// NewConnector creates new Connector
// gateways is a map of preconfigured gateway pool
// error is a chan of generated errors by Connector
func NewConnector(gateways *[]model.Gateway, egress *chan model.Message, error *chan error) *Connector {
	return &Connector{
		gateways: gateways,
		egress:   egress,
		timer:    time.NewTimer(time.Duration(time.Second)),
		Error:    error,
	}
}

// GetConnection takes connection from pool by gateway name
func (c *Connector) GetConnection(gateway string) (*connection.Connection, error) {
	if conn, exists := c.pool[gateway]; exists {
		return conn, nil
	}
	return nil, errors.New(fmt.Sprintf("No %s connection found ", gateway))
}

// StartConnection starts connection by gateway name
func (c *Connector) StartConnection(gateway *model.Gateway) error {
	stop := make(chan int, 1)
	inbox := make(chan model.Message, gateway.InboxSize)
	conn, err := connection.NewConnection(gateway, &inbox, c.egress, &stop, c.Error)
	if err != nil {
		return err
	}
	c.pool[gateway.Name] = conn
	go c.pool[gateway.Name].Run()
	return nil
}

// StopConnection stops connection by gateway name
func (c *Connector) StopConnection(gateway string) error {
	_, err := c.GetConnection(gateway)
	if err != nil {
		return err
	}
	delete(c.pool, gateway)
	return nil
}

// Run starts a pool pool
func (c *Connector) Run() {
	for _, gateway := range *c.gateways {
		err := c.StartConnection(&gateway)
		if err != nil {
			*c.Error <- err
		}
	}
	for {
		for gateway, conn := range c.pool {
			select {
			case err := <-*conn.Error:
				log.Printf("Error ocured in gateway %s connection: %s", gateway, err.Error())
				delete(c.pool, gateway)
				log.Printf("Gateway %s removed from pool.", gateway)
			case <-c.timer.C:
				c.timer.Reset(time.Duration(time.Second))
			}
		}
	}
}
