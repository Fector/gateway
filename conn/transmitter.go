package conn

import (
	"github.com/DeathHand/smpp/pdu"
	"github.com/DeathHand/smpp/protocol"
	"net"
)

type Transmitter struct {
	writer *protocol.Writer
	t      *chan pdu.Pdu
}

func NewTransmitter(c *net.TCPConn, t *chan pdu.Pdu) *Transmitter {
	return &Transmitter{
		writer: protocol.NewWriter(c),
		t:      t,
	}
}

func (t *Transmitter) Transmit() {

}
