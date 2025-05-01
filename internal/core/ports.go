package core

import "context"

type Storage interface {
	GetAllMessages(ctx context.Context) ([]Message, error)
	SaveMessage(ctx context.Context, message Message) error
}

type Scheduler interface {
	ScheduleMessage(ctx context.Context, message Message) error
	RetryMessage(ctx context.Context, id string) error
}
