package conn

import (
	"github.com/DeathHand/smpp/io"
	"github.com/DeathHand/smpp/pdu"
	"net"
)

type Transmitter struct {
	Writer  *io.Writer
	Ingress *chan pdu.Pdu
}

func NewTransmitter(c *net.TCPConn, i *chan pdu.Pdu) *Transmitter {
	return &Transmitter{
		Writer:  io.NewWriter(c),
		Ingress: i,
	}
}

func (t *Transmitter) Transmit() {

}
