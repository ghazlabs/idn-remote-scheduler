package core

import (
	"context"
	"fmt"

	"gopkg.in/validator.v2"
)

type Service interface {
	GetAllMessages(ctx context.Context, msg Message) ([]Message, error)
	SendMessage(ctx context.Context, msg Message) error
	RetryMessage(ctx context.Context, msg Message) error
}

type ServiceConfig struct {
	Storage   Storage   `validate:"nonnil"`
	Scheduler Scheduler `validate:"nonnil"`
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

func (s *service) GetAllMessages(ctx context.Context, msg Message) ([]Message, error) {
	messages, err := s.Storage.GetAllMessages(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	return messages, nil
}

func (s *service) SendMessage(ctx context.Context, msg Message) error {
	err := s.Storage.SaveMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to schedule message: %w", err)
	}

	return nil
}

func (s *service) RetryMessage(ctx context.Context, msg Message) error {
	err := s.Storage.UpdateMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}
