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
	rx             *chan []byte
	reader         *protocol.Reader
	writer         *protocol.Writer
	gateway        *model.Gateway
	Ingress        *chan model.Message
	Egress         *chan model.Message
	timer          *time.Ticker
	logger         *log.Logger
	Stop           *chan int
	Error          *chan error
}

func NewConnection(
	gateway *model.Gateway,
	rx *chan []byte,
	ingress *chan model.Message,
	egress *chan model.Message,
	stop *chan int,
	error *chan error,
) (*Connection, error) {
	return &Connection{
		gateway: gateway,
		rx:      rx,
		Ingress: ingress,
		Egress:  egress,
		timer:   time.NewTicker(time.Duration(gateway.EnquireLinkTime) * time.Second),
		logger: log.New(
			os.Stdout,
			fmt.Sprintf("%s: ", gateway.Name),
			log.Ldate|log.Ltime|log.Lmicroseconds,
		),
		Stop:  stop,
		Error: error,
	}, nil
}

func (c *Connection) nextSequence() uint32 {
	return atomic.AddUint32(c.SequenceNumber, 1)
}

func (c *Connection) updateTime() {
	atomic.StoreInt64(c.LastPacketTime, time.Now().UnixNano())
}

func (c *Connection) connect() error {
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
	c.reader = protocol.NewReader(tcpConn)
	c.writer = protocol.NewWriter(tcpConn)
	c.logger.Printf("Connecting to %s successful", addr)
	return nil
}

func (c *Connection) getBindPdu() (pdu.Pdu, error) {
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
		return pdu.NewBindReceiver(c.nextSequence(), &body), nil
	case protocol.BindModeTX:
		return pdu.NewBindTransmitter(c.nextSequence(), &body), nil
	case protocol.BindModeTRX:
		return pdu.NewBindTransceiver(c.nextSequence(), &body), nil
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

func (c *Connection) bind() error {
	c.logger.Printf("Binding connection")
	req, err := c.getBindPdu()
	if err != nil {
		return err
	}
	err = c.SetDeadline(time.Now().Add(time.Duration(c.gateway.BindTimeout) * time.Second))
	if err != nil {
		return err
	}
	_, err = c.writer.WritePdu(&req)
	if err != nil {
		return err
	}
	err = c.SetDeadline(time.Now().Add(time.Duration(c.gateway.BindTimeout) * time.Second))
	if err != nil {
		return err
	}
	p, err := c.reader.ReadPacket()
	if err != nil {
		return err
	}
	resp, err := c.reader.ReadPdu(p)
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

func (c *Connection) unbind() error {
	return nil
}

func (c *Connection) receive() {
	for {
		data, err := c.reader.ReadPacket()
		if err != nil {
			*c.Error <- err
			continue
		}
		*c.rx <- *data
	}
}

func (c Connection) handlePacket(p *[]byte) {

}

func (c *Connection) handleMessage(m *model.Message) {

}

func (c *Connection) enquireLink() error {
	req := pdu.NewEnquireLink(c.nextSequence())
	_, err := c.writer.WritePdu(&req)
	if err != nil {
		return err
	}
	err = c.TCPConn.SetReadDeadline(time.Now().Add(time.Duration(c.gateway.BindTimeout) * time.Second))
	if err != nil {
		return err
	}
	data, err := c.reader.ReadPacket()
	if err != nil {
		return err
	}
	rep, err := c.reader.ReadPdu(data)
	if err != nil {
		return err
	}
	if response, ok := rep.(pdu.EnquireLinkResp); ok {
		if response.Header.CommandStatus != protocol.EsmeRok {
			return errors.New(
				fmt.Sprintf("Wrong link response status: %s", protocol.GetStatusName(response.Header.CommandStatus)),
			)
		}
		return nil
	}
	return errors.New("Enquire Link failed. Wrong response operation ")
}

func (c *Connection) Run() {
	err := c.connect()
	if err != nil {
		*c.Error <- err
		return
	}
	err = c.bind()
	if err != nil {
		*c.Error <- err
		return
	}
	go c.receive()
	for {
		select {
		case m := <-*c.Ingress:
			go c.handleMessage(&m)
		case p := <-*c.rx:
			go c.handlePacket(&p)
		case <-c.timer.C:
			err := c.enquireLink()
			if err != nil {
				*c.Error <- err
				return
			}
		case <-*c.Stop:
			c.logger.Printf("Stop flag received. Unbinding.")
			err := c.unbind()
			if err != nil {
				*c.Error <- err
				return
			}
			c.logger.Printf("Connection unbond.")
			err = c.TCPConn.Close()
			if err != nil {
				*c.Error <- err
				return
			}
			c.logger.Printf("Connection closed.")
			return
		}
	}
}
