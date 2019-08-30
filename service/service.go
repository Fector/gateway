package service

import (
	"github.com/DeathHand/gateway/api"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/gateway/router"
)

type Service struct {
	Ingress *chan model.Message
	Egress  *chan model.Message
	Memory  memory.Memory
	ErrChan *chan error
}

func New(c *Config) (*Service, error) {
	ingress := make(chan model.Message, c.IngressSize)
	egress := make(chan model.Message, c.EgressSize)
	errChan := make(chan error, 1)
	m, err := memory.NewMapMemory("/var/lib/memory", &errChan)
	if err != nil {
		return nil, err
	}
	return &Service{
		Ingress: &ingress,
		Egress:  &egress,
		Memory:  m,
		ErrChan: &errChan,
	}, nil
}

func (s *Service) Run() {
	go s.Memory.Run()
	go router.NewRouter(s).Run()
	go api.NewApi(s).Serve()
}
