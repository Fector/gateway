package connection

import (
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/smpp/pdu"
	"github.com/DeathHand/smpp/protocol"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"
)

type Connection struct {
	*net.TCPConn
	SequenceNumber *uint32
	LastPacketTime *int64
	rx             *protocol.Reader
	tx             *protocol.Writer
	gateway        *model.Gateway
	inbox          *chan model.Message
	outbox         *chan model.Message
	timer          *time.Ticker
	stop           *chan int
	error          *chan error
	logger         *log.Logger
}

func NewConnection(
	gateway *model.Gateway,
	inbox *chan model.Message,
	outbox *chan model.Message,
	stop *chan int,
	error *chan error,
) (*Connection, error) {
	return &Connection{
		gateway: gateway,
		inbox:   inbox,
		outbox:  outbox,
		timer:   time.NewTicker(time.Duration(gateway.EnquireLinkTime) * time.Second),
		stop:    stop,
		error:   error,
		logger: log.New(
			os.Stdout,
			fmt.Sprintf("%s: ", gateway.Name),
			log.Ldate|log.Ltime|log.Lmicroseconds,
		),
	}, nil
}

func (c *Connection) Gateway() *model.Gateway {
	return c.gateway
}

func (c *Connection) SendMessage(message *model.Message) {
	*c.inbox <- *message
}

func (c *Connection) NextSequence() uint32 {
	return atomic.AddUint32(c.SequenceNumber, 1)
}

func (c *Connection) UpdateTime() {
	atomic.StoreInt64(c.LastPacketTime, time.Now().UnixNano())
}

func (c *Connection) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.gateway.Host, c.gateway.Port)
	c.logger.Printf("Connecting to %s", addr)
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
	c.rx = protocol.NewReader(tcpConn)
	c.tx = protocol.NewWriter(tcpConn)
	c.logger.Printf("Connecting to %s successful", addr)
	return nil
}

func (c *Connection) Stop() {
	*c.stop <- 1
}

func (c *Connection) getBindPdu() (pdu.Pdu, error) {
	header := pdu.Header{
		CommandStatus:  protocol.EsmeRok,
		SequenceNumber: c.NextSequence(),
	}
	body := pdu.BindBody{
		SystemId:         c.gateway.SystemId,
		Password:         c.gateway.Password,
		SystemType:       c.gateway.SystemType,
		InterfaceVersion: c.gateway.InterfaceVersion,
		AddrTon:          c.gateway.AddrTon,
		AddrNpi:          c.gateway.AddrNpi,
		AddressRange:     c.gateway.AddressRange,
	}
	switch c.gateway.BindMode {
	case protocol.BindModeRX:
		header.CommandId = protocol.BindReceiver
		return &pdu.BindReceiver{
			Header: &header,
			Body:   &body,
			Tlv:    nil,
		}, nil
	case protocol.BindModeTX:
		header.CommandId = protocol.BindTransmitter
		return &pdu.BindTransmitter{
			Header: &header,
			Body:   &body,
			Tlv:    nil,
		}, nil
	case protocol.BindModeTRX:
		header.CommandId = protocol.BindTransceiver
		return &pdu.BindTransceiver{
			Header: &header,
			Body:   &body,
			Tlv:    nil,
		}, nil
	}
	return nil, errors.New("Unknown bind mode ")
}

func (c *Connection) checkBindResp(header *pdu.Header) error {
	if header.CommandStatus != protocol.EsmeRok {
		return errors.New(
			fmt.Sprintf("Bind failed with code: %s", protocol.GetStatusName(header.CommandStatus)),
		)
	}
	return nil
}

func (c *Connection) Bind() error {
	c.logger.Printf("Binding connection")
	req, err := c.getBindPdu()
	if err != nil {
		return err
	}
	err = c.SetDeadline(time.Now().Add(time.Duration(c.gateway.BindTimeout) * time.Second))
	if err != nil {
		return err
	}
	_, err = c.tx.WritePdu(&req)
	if err != nil {
		return err
	}
	err = c.SetDeadline(time.Now().Add(time.Duration(c.gateway.BindTimeout) * time.Second))
	if err != nil {
		return err
	}
	p, err := c.rx.ReadPacket()
	if err != nil {
		return err
	}
	resp, err := c.rx.ReadPdu(p)
	switch p := resp.(type) {
	case pdu.BindReceiverResp:
		return c.checkBindResp(p.Header)
	case pdu.BindTransceiverResp:
		return c.checkBindResp(p.Header)
	case pdu.BindTransmitterResp:
		return c.checkBindResp(p.Header)
	}
	if err != nil {
		return err
	}
	c.logger.Printf("Connection bond successful")
	return nil
}

func (c *Connection) handleMessage(m *model.Message) {

}

func (c *Connection) enquireLink() {

}

func (c *Connection) Run() {
	for {
		select {
		case m := <-*c.inbox:
			go c.handleMessage(&m)
		case <-c.timer.C:
			go c.enquireLink()
		case <-*c.stop:
			*c.error <- errors.New("Stop flag received ")
		}
	}
}
