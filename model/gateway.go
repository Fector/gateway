package model

type Gateway struct {
	Name               string
	Host               string
	Port               int
	BindMode           string
	BindTimeout        int
	SystemId           string
	Password           string
	SystemType         string
	InterfaceVersion   uint32
	AddrTon            uint32
	AddrNpi            uint32
	AddressRange       string
	ReadChanSize       int
	EnquireLinkTime    int
	EnquireLinkTimeout int
	InboxSize          int
	OutboxSize         int
}
