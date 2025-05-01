package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ghazlabs/idn-remote-scheduler/internal/core"
	"github.com/go-co-op/gocron/v2"
	"gopkg.in/validator.v2"
)

const (
	// DefaultMaxRetries is the default number of retries for a message
	DefaultMaxRetries = 3
	// DefaultRetryDelay is the default delay between retries
	DefaultRetryDelay = 30 * time.Second
)

type GoCronScheduler struct {
	GoCronSchedulerConfig
}

type GoCronSchedulerConfig struct {
	Client    gocron.Scheduler `validate:"nonnil"`
	Publisher Publisher        `validate:"nonnil"`
	Storage   core.Storage     `validate:"nonnil"`
}

func NewGoCronScheduler(cfg GoCronSchedulerConfig) (*GoCronScheduler, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &GoCronScheduler{
		GoCronSchedulerConfig: cfg,
	}, nil
}

func (s *GoCronScheduler) ScheduleMessage(ctx context.Context, msg core.Message, at int64) error {
	isNow := time.Now().Unix() >= at

	var option gocron.OneTimeJobStartAtOption
	if isNow {
		option = gocron.OneTimeJobStartImmediately()
	} else {
		option = gocron.OneTimeJobStartDateTime(time.Unix(at, 0))
	}

	// schedule message
	_, err := s.Client.NewJob(
		gocron.OneTimeJob(option),
		gocron.NewTask(func() {
			s.sendMessage(context.Background(), msg)
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to schedule message: %w", err)
	}

	return nil
}

func (s *GoCronScheduler) sendMessage(ctx context.Context, msg core.Message) {
	now := time.Now().Unix()
	err := s.Publisher.Publish(context.Background(), msg)
	if err != nil {
		// if message failed to send more than max retries
		// or if the error is session expired
		// mark the message as failed
		if (msg.RetriedCount >= DefaultMaxRetries) || (err == core.ErrSessionExpired) {
			reason := fmt.Sprintf("failed to send message after %d retries: %v", DefaultMaxRetries, err)
			msg.Reason = &reason
			msg.Status = core.MessageStatusFailed

			slog.Error("failed to send message", slog.String("message", msg.String()), slog.String("err", err.Error()))
		} else {
			// otherwise, retry the message
			s.retryMessage(ctx, msg)
			return
		}
	} else {
		msg.Status = core.MessageStatusSent
		msg.SentAt = &now
	}

	err = s.Storage.UpdateMessage(ctx, msg)
	if err != nil {
		slog.Error("failed to update message", slog.String("message", msg.String()), slog.String("err", err.Error()))
		return
	}
}

func (s *GoCronScheduler) retryMessage(ctx context.Context, msg core.Message) error {
	scheduleTime := time.Now().Add(DefaultRetryDelay).Unix()
	err := s.ScheduleMessage(ctx, msg, scheduleTime)
	if err != nil {
		return fmt.Errorf("failed to retry message: %w", err)
	}

	msg.RetriedCount++
	msg.ScheduledSendingAt = scheduleTime
	msg.Status = core.MessageStatusScheduled
	err = s.Storage.UpdateMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}
