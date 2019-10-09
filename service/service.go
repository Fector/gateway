package service

import (
	"github.com/DeathHand/gateway/api"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/gateway/router"
)

type Service struct {
	Config  *Config
	Ingress *chan model.Message
	Egress  *chan model.Message
	Memory  memory.Memory
	Error   *chan error
}

func New(c *Config) (*Service, error) {
	errors := make(chan error, 1)
	ingress := make(chan model.Message, c.IngressSize)
	egress := make(chan model.Message, c.EgressSize)
	m, err := memory.NewMapMemory(c.MapMemory.Path, &errors)
	if err != nil {
		return nil, err
	}
	return &Service{
		Config:  c,
		Ingress: &ingress,
		Egress:  &egress,
		Memory:  m,
		Error:   &errors,
	}, nil
}

func (s *Service) Run() {
	go s.Memory.Run()
	go api.NewApi(s.Ingress, s.Memory, s.Error).Run()
	go router.NewRouter(&s.Config.Gateways, s.Ingress, s.Egress, s.Memory, s.Error).Run()
}
