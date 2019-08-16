package model

import "github.com/DeathHand/smpp/pdu"

type Message struct {
	Uuid      string
	Body      *pdu.Pdu
	Ttl       int
	Timestamp int64
}
