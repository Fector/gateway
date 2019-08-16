package router

type Driver interface {
	Init() error
	Run() error
	CallBack() error
	Close() error
}
