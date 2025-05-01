package core

import (
	"context"
	"fmt"

	"gopkg.in/validator.v2"
)

type Service interface {
	GetAllMessages(ctx context.Context, req GetAllMessageRequest) ([]Message, error)
	SendMessage(ctx context.Context, req SendMessageRequest) error
	RetryMessage(ctx context.Context, req RetryMessageRequest) error
}

type ServiceConfig struct {
	Storage Storage `validate:"nonnil"`
}

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

func (s *service) GetAllMessages(ctx context.Context, req GetAllMessageRequest) ([]Message, error) {
	messages, err := s.Storage.GetAllMessages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	return messages, nil
}

func (s *service) SendMessage(ctx context.Context, req SendMessageRequest) error {
	return nil
}

func (s *service) RetryMessage(ctx context.Context, req RetryMessageRequest) error {
	return nil
}
