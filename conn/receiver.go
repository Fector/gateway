package conn

import (
	"github.com/DeathHand/smpp/pdu"
	"github.com/DeathHand/smpp/protocol"
	"net"
)

type Receiver struct {
	reader *protocol.Reader
	egress *chan pdu.Pdu
}

func NewReceiver(c *net.TCPConn, e *chan pdu.Pdu) *Receiver {
	return &Receiver{
		reader: protocol.NewReader(c),
		egress: e,
	}
}

func (r *Receiver) Receive() {

}
