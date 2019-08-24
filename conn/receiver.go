package conn

import (
	"github.com/DeathHand/smpp/pdu"
	"github.com/DeathHand/smpp/protocol"
	"net"
)

type Receiver struct {
	reader *protocol.Reader
	r      *chan pdu.Pdu
}

func NewReceiver(c *net.TCPConn, r *chan pdu.Pdu) *Receiver {
	return &Receiver{
		reader: protocol.NewReader(c),
		r:      r,
	}
}

func (r *Receiver) Receive() {

}
