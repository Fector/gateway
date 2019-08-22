package conn

import (
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/smpp/pdu"
	"net"
	"time"
)

type Connection struct {
	rx      *Receiver
	tx      *Transmitter
	gateway *model.Gateway
	context *Context
	logger  *Logger
	inbox   *chan model.Message
	timer   *time.Ticker
	stop    *chan int
	error   *chan error
}

func NewConnection(gateway *model.Gateway, inbox *chan model.Message, stop *chan int) (*Connection, error) {
	return &Connection{
		gateway: gateway,
		context: &Context{},
		logger:  NewLogger(gateway.Name),
		inbox:   inbox,
		timer:   time.NewTicker(time.Duration(gateway.EnquireLinkTime) * time.Second),
		stop:    stop,
	}, nil
}

func (c *Connection) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.gateway.Host, c.gateway.Port)
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		return err
	}
	tcpConn, isTcp := conn.(*net.TCPConn)
	if !isTcp {
		return errors.New("Unknown connection type ")
	}
	ingress := make(chan pdu.Pdu, c.gateway.IngressSize)
	egress := make(chan pdu.Pdu, c.gateway.EgressSize)
	c.rx = NewReceiver(tcpConn, &egress)
	c.tx = NewTransmitter(tcpConn, &ingress)
	return nil
}

func (c *Connection) Bind() error {
	return NewAuth(c).Auth()
}

func (c *Connection) handlePdu(p *pdu.Pdu) {
	*c.tx.ingress <- p
}

func (c *Connection) handleMessage(m *model.Message) {

}

func (c *Connection) enquireLink() {

}

func (c *Connection) Run() error {
	go c.rx.Receive()
	go c.tx.Transmit()
	for {
		select {
		case err := <-*c.error:
			return err
		case p := <-*c.rx.egress:
			go c.handlePdu(&p)
		case m := <-*c.inbox:
			go c.handleMessage(&m)
		case <-c.timer.C:
			go c.enquireLink()
		case <-*c.stop:
			return errors.New("Stop flag received ")
		}
	}
}
