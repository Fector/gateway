package service

import "github.com/DeathHand/gateway/storage"

type Service struct {
	GatewayStorage *storage.Gateway
	MessageStorage *storage.Message
}

func NewService() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) Run() error {
	go s.MessageStorage.Collect()

	return nil
}
