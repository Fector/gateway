package service

import (
	"github.com/DeathHand/gateway/api"
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"github.com/DeathHand/gateway/router"
)

type Service struct {
	config *Config
	error  *chan error
}

func (s *Service) Run() error {
	ingress := map[string]*chan model.Message{}
	egress := map[string]*chan model.Message{}
	mtIngress := make(chan model.Message, s.config.MtIngressSize)
	moIngress := make(chan model.Message, s.config.MoIngressSize)
	mtEgress := make(chan model.Message, s.config.MtEgressSize)
	moEgress := make(chan model.Message, s.config.MoEgressSize)
	ingress["mt"] = &mtIngress
	ingress["mo"] = &moIngress
	egress["mt"] = &mtEgress
	egress["mo"] = &moEgress
	m, err := memory.NewMapMemory("/var/lib/memory")
	if err != nil {
		return err
	}
	go router.NewRouter(ingress, egress, m).Run()
	go api.NewApi(ingress["mt"], m).Serve()
	return nil
}
