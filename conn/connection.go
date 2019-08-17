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
	Rx      *Receiver
	Tx      *Transmitter
	Gateway *model.Gateway
	Context *Context
	Logger  *Logger
	Inbox   *chan model.Message
	Stop    *chan int
	Error   *chan error
}

func NewConnection(gateway *model.Gateway, inbox *chan model.Message, stop *chan int) (*Connection, error) {
	return &Connection{
		Gateway: gateway,
		Context: &Context{},
		Logger:  NewLogger(gateway.Name),
		Inbox:   inbox,
		Stop:    stop,
	}, nil
}

func (c *Connection) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.Gateway.Host, c.Gateway.Port)

	d := net.Dialer{Timeout: 5 * time.Second}

	conn, err := d.Dial("tcp", addr)

	if err != nil {
		return err
	}

	tcpConn, isTcp := conn.(*net.TCPConn)

	if !isTcp {
		return errors.New("Unknown connection type ")
	}

	ingress := make(chan pdu.Pdu, c.Gateway.IngressSize)
	egress := make(chan pdu.Pdu, c.Gateway.EgressSize)

	c.Rx = NewReceiver(tcpConn, &egress)
	c.Tx = NewTransmitter(tcpConn, &ingress)

	return nil
}

func (c *Connection) Bind() error {
	return NewAuth(c).Auth()
}

func (c *Connection) handlePdu(p *pdu.Pdu) {
	*c.Tx.Ingress <- p
}

func (c *Connection) handleMessage(m *model.Message) {

}

func (c *Connection) enquireLink() {

}

func (c *Connection) Run() error {
	go c.Rx.Receive()
	go c.Tx.Transmit()

	linkTimer := time.NewTicker(time.Duration(c.Gateway.EnquireLinkTime) * time.Second)

	for {
		select {
		case err := <-*c.Error:
			return err
		case p := <-*c.Rx.Egress:
			c.handlePdu(&p)
		case m := <-*c.Inbox:
			c.handleMessage(&m)
		case <-linkTimer.C:
			c.enquireLink()
		case <-*c.Stop:
			return errors.New("Stop flag received ")
		}
	}
}
