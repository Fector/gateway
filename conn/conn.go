package conn

import (
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/smpp/pdu"
	"net"
	"time"
)

type Conn struct {
	*net.TCPConn
	rx      *Receiver
	tx      *Transmitter
	gateway *model.Gateway
	context *Context
	inbox   *chan model.Message
	outbox  *chan model.Message
	timer   *time.Ticker
	stop    *chan int
	error   *chan error
}

func NewConnection(gateway *model.Gateway, inbox *chan model.Message, outbox *chan model.Message, stop *chan int) (*Conn, error) {
	return &Conn{
		gateway: gateway,
		context: &Context{},
		inbox:   inbox,
		outbox:  outbox,
		timer:   time.NewTicker(time.Duration(gateway.EnquireLinkTime) * time.Second),
		stop:    stop,
	}, nil
}

func (c *Conn) Connect() error {
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
	c.TCPConn = tcpConn
	r := make(chan pdu.Pdu, c.gateway.ReadChanSize)
	t := make(chan pdu.Pdu, c.gateway.WriteChanSize)
	c.rx = NewReceiver(tcpConn, &r)
	c.tx = NewTransmitter(tcpConn, &t)
	return nil
}

func (c *Conn) handlePdu(p *pdu.Pdu) {
	*c.tx.t <- p
}

func (c *Conn) handleMessage(m *model.Message) {

}

func (c *Conn) enquireLink() {

}

func (c *Conn) Run() error {
	go c.rx.Receive()
	go c.tx.Transmit()
	for {
		select {
		case err := <-*c.error:
			return err
		case p := <-*c.rx.r:
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
