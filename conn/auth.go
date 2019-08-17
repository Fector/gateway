package conn

import (
	"errors"
	"fmt"
	"github.com/DeathHand/smpp/pdu"
	"github.com/DeathHand/smpp/protocol"
	"time"
)

type Auth struct {
	conn *Connection
}

func NewAuth(c *Connection) *Auth {
	return &Auth{conn: c}
}

func (a *Auth) getBindPDU() (pdu.Pdu, error) {
	switch a.conn.Gateway.BindMode {
	case protocol.BindModeRX:
		return &pdu.BindReceiver{
			Header: &pdu.Header{
				CommandId:      protocol.BindReceiver,
				CommandStatus:  protocol.EsmeRok,
				SequenceNumber: a.conn.Context.NextSequence(),
			},
			Body: &pdu.BindBody{
				SystemId:         a.conn.Gateway.SystemId,
				Password:         a.conn.Gateway.Password,
				SystemType:       a.conn.Gateway.SystemType,
				InterfaceVersion: a.conn.Gateway.InterfaceVersion,
				AddrTon:          a.conn.Gateway.AddrTon,
				AddrNpi:          a.conn.Gateway.AddrNpi,
				AddressRange:     a.conn.Gateway.AddressRange,
			},
			Tlv: nil,
		}, nil
	case protocol.BindModeTX:
		return &pdu.BindTransmitter{
			Header: &pdu.Header{
				CommandId:      protocol.BindTransmitter,
				CommandStatus:  protocol.EsmeRok,
				SequenceNumber: a.conn.Context.NextSequence(),
			},
			Body: &pdu.BindBody{
				SystemId:         a.conn.Gateway.SystemId,
				Password:         a.conn.Gateway.Password,
				SystemType:       a.conn.Gateway.SystemType,
				InterfaceVersion: a.conn.Gateway.InterfaceVersion,
				AddrTon:          a.conn.Gateway.AddrTon,
				AddrNpi:          a.conn.Gateway.AddrNpi,
				AddressRange:     a.conn.Gateway.AddressRange,
			},
			Tlv: nil,
		}, nil
	case protocol.BindModeTRX:
		return &pdu.BindTransceiver{
			Header: &pdu.Header{
				CommandId:      protocol.BindTransceiver,
				CommandStatus:  protocol.EsmeRok,
				SequenceNumber: a.conn.Context.NextSequence(),
			},
			Body: &pdu.BindBody{
				SystemId:         a.conn.Gateway.SystemId,
				Password:         a.conn.Gateway.Password,
				SystemType:       a.conn.Gateway.SystemType,
				InterfaceVersion: a.conn.Gateway.InterfaceVersion,
				AddrTon:          a.conn.Gateway.AddrTon,
				AddrNpi:          a.conn.Gateway.AddrNpi,
				AddressRange:     a.conn.Gateway.AddressRange,
			},
			Tlv: nil,
		}, nil
	}

	return nil, errors.New("Unknown bind mode ")
}

func (a *Auth) checkBindResp(header *pdu.Header) error {
	if header.CommandStatus != protocol.EsmeRok {
		return errors.New(
			fmt.Sprintf("Bind failed with code: %s", protocol.GetStatusName(header.CommandStatus)),
		)
	}

	return nil
}

func (a *Auth) Auth() error {
	req, err := a.getBindPDU()

	if err != nil {
		return err
	}

	err = a.conn.Rx.Reader.SetDeadline(time.Now().Add(time.Duration(1 * time.Second)))

	if err != nil {
		return err
	}

	_, err = a.conn.Tx.Writer.WritePdu(&req)

	if err != nil {
		return err
	}

	err = a.conn.Rx.Reader.SetDeadline(time.Now().Add(time.Duration(1 * time.Second)))

	if err != nil {
		return err
	}

	resp, err := a.conn.Rx.Reader.ReadPdu()

	switch p := resp.(type) {
	case pdu.BindReceiverResp:
		return a.checkBindResp(p.Header)
	case pdu.BindTransceiverResp:
		return a.checkBindResp(p.Header)
	case pdu.BindTransmitterResp:
		return a.checkBindResp(p.Header)
	}

	if err != nil {
		return err
	}

	return nil
}
