package model

type Gateway struct {
	Uuid             string
	Name             string
	Host             string
	Port             int
	BindMode         string
	SystemId         string
	Password         string
	SystemType       string
	InterfaceVersion uint32
	AddrTon          uint32
	AddrNpi          uint32
	AddressRange     string
	IngressSize      int
	EgressSize       int
	EnquireLinkTime  int
}
