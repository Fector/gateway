package conn

import (
	"errors"
	"fmt"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/smpp/pdu"
	"github.com/DeathHand/smpp/protocol"
	"log"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
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
}

func NewConnection(
	gateway *model.Gateway,
	inbox *chan model.Message,
	outbox *chan model.Message,
	stop *chan int,
	error *chan error,
) (*Conn, error) {
	return &Conn{
		gateway: gateway,
		inbox:   inbox,
		outbox:  outbox,
		timer:   time.NewTicker(time.Duration(gateway.EnquireLinkTime) * time.Second),
		stop:    stop,
		error:   error,
	}, nil
}

func (c *Conn) NextSequence() uint32 {
	return atomic.AddUint32(c.SequenceNumber, 1)
}

func (c *Conn) UpdateTime() {
	atomic.StoreInt64(c.LastPacketTime, time.Now().UnixNano())
}

func (c *Conn) Log(v ...interface{}) {
	log.Println(fmt.Sprintf("[gateway:%s] ", c.gateway.Name), v)
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
	c.rx = protocol.NewReader(tcpConn)
	c.tx = protocol.NewWriter(tcpConn)
	return nil
}

func (c *Conn) getBindPdu() (pdu.Pdu, error) {
	switch c.gateway.BindMode {
	case protocol.BindModeRX:
		return &pdu.BindReceiver{
			Header: &pdu.Header{
				CommandId:      protocol.BindReceiver,
				CommandStatus:  protocol.EsmeRok,
				SequenceNumber: c.NextSequence(),
			},
			Body: &pdu.BindBody{
				SystemId:         c.gateway.SystemId,
				Password:         c.gateway.Password,
				SystemType:       c.gateway.SystemType,
				InterfaceVersion: c.gateway.InterfaceVersion,
				AddrTon:          c.gateway.AddrTon,
				AddrNpi:          c.gateway.AddrNpi,
				AddressRange:     c.gateway.AddressRange,
			},
			Tlv: nil,
		}, nil
	case protocol.BindModeTX:
		return &pdu.BindTransmitter{
			Header: &pdu.Header{
				CommandId:      protocol.BindTransmitter,
				CommandStatus:  protocol.EsmeRok,
				SequenceNumber: c.NextSequence(),
			},
			Body: &pdu.BindBody{
				SystemId:         c.gateway.SystemId,
				Password:         c.gateway.Password,
				SystemType:       c.gateway.SystemType,
				InterfaceVersion: c.gateway.InterfaceVersion,
				AddrTon:          c.gateway.AddrTon,
				AddrNpi:          c.gateway.AddrNpi,
				AddressRange:     c.gateway.AddressRange,
			},
			Tlv: nil,
		}, nil
	case protocol.BindModeTRX:
		return &pdu.BindTransceiver{
			Header: &pdu.Header{
				CommandId:      protocol.BindTransceiver,
				CommandStatus:  protocol.EsmeRok,
				SequenceNumber: c.NextSequence(),
			},
			Body: &pdu.BindBody{
				SystemId:         c.gateway.SystemId,
				Password:         c.gateway.Password,
				SystemType:       c.gateway.SystemType,
				InterfaceVersion: c.gateway.InterfaceVersion,
				AddrTon:          c.gateway.AddrTon,
				AddrNpi:          c.gateway.AddrNpi,
				AddressRange:     c.gateway.AddressRange,
			},
			Tlv: nil,
		}, nil
	}
	return nil, errors.New("Unknown bind mode ")
}

func (c *Conn) checkBindResp(header *pdu.Header) error {
	if header.CommandStatus != protocol.EsmeRok {
		return errors.New(
			fmt.Sprintf("Bind failed with code: %s", protocol.GetStatusName(header.CommandStatus)),
		)
	}
	return nil
}

func (c *Conn) Bind() error {
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
	return nil
}

func (c *Conn) handleMessage(m *model.Message) {

}

func (c *Conn) enquireLink() {

}

func (c *Conn) Run() error {
	for {
		select {
		case m := <-*c.inbox:
			go c.handleMessage(&m)
		case <-c.timer.C:
			go c.enquireLink()
		case <-*c.stop:
			return errors.New("Stop flag received ")
		}
	}
}
