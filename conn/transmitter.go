package conn

import (
	"github.com/DeathHand/smpp/pdu"
	"github.com/DeathHand/smpp/protocol"
	"net"
)

type Transmitter struct {
	writer  *protocol.Writer
	ingress *chan pdu.Pdu
}

func NewTransmitter(c *net.TCPConn, i *chan pdu.Pdu) *Transmitter {
	return &Transmitter{
		writer:  protocol.NewWriter(c),
		ingress: i,
	}
}

func (t *Transmitter) Transmit() {

}
