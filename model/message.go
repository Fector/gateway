package model

import "github.com/DeathHand/smpp/pdu"

type Message struct {
	Uuid      string
	Gateway   string
	Type      string
	Body      *pdu.Pdu
	Ttl       int64
	Timestamp int64
	Callback  string
}
