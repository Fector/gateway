package engine

type Router struct {
	Ingress *chan interface{}
	Egress  *chan interface{}
}
