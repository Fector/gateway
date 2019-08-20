package conn

import (
	"github.com/DeathHand/smpp/io"
	"github.com/DeathHand/smpp/pdu"
	"net"
)

type Receiver struct {
	reader *io.Reader
	egress *chan pdu.Pdu
}

func NewReceiver(c *net.TCPConn, e *chan pdu.Pdu) *Receiver {
	return &Receiver{
		reader: io.NewReader(c),
		egress: e,
	}
}

func (r *Receiver) Receive() {

}
