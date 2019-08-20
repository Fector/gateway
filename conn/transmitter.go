package conn

import (
	"github.com/DeathHand/smpp/io"
	"github.com/DeathHand/smpp/pdu"
	"net"
)

type Transmitter struct {
	writer  *io.Writer
	ingress *chan pdu.Pdu
}

func NewTransmitter(c *net.TCPConn, i *chan pdu.Pdu) *Transmitter {
	return &Transmitter{
		writer:  io.NewWriter(c),
		ingress: i,
	}
}

func (t *Transmitter) Transmit() {

}
