package core

import (
	"context"

	"gopkg.in/validator.v2"
)

type Service interface {
	GetAllMessages(ctx context.Context) ([]Message, error)
	SendMessage(ctx context.Context, id string) error
	RetryMessage(ctx context.Context, id string) error
}

type ServiceConfig struct{}

type service struct {
	ServiceConfig
}

func NewService(config ServiceConfig) (Service, error) {
	err := validator.Validate(config)
	if err != nil {
		return nil, err
	}

	return &service{
		ServiceConfig: config,
	}, nil
}

func (s *service) GetAllMessages(ctx context.Context) ([]Message, error) {
	return nil, nil
}

func (s *service) SendMessage(ctx context.Context, id string) error {
	return nil
}

func (s *service) RetryMessage(ctx context.Context, id string) error {
	return nil
}
