package conn

import (
	"fmt"
	"github.com/DeathHand/gateway/model"
	"net"
)

type Connection struct {
	*net.TCPConn
	Gateway *model.Gateway
	Context *Context
}

func NewConnection(gateway *model.Gateway) (*Connection, error) {
	addr, err := net.ResolveTCPAddr(
		"tcp",
		fmt.Sprintf("%s:%d", gateway.Host, gateway.Port),
	)

	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return nil, err
	}

	return &Connection{
		TCPConn: conn,
		Gateway: gateway,
		Context: &Context{},
	}, nil
}
