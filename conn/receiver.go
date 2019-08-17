package conn

import (
	"github.com/DeathHand/smpp/io"
	"github.com/DeathHand/smpp/pdu"
	"net"
)

type Receiver struct {
	Reader *io.Reader
	Egress *chan pdu.Pdu
}

func NewReceiver(c *net.TCPConn, e *chan pdu.Pdu) *Receiver {
	return &Receiver{
		Reader: io.NewReader(c),
		Egress: e,
	}
}

func (r *Receiver) Receive() {

}
