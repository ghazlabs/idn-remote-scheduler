package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ghazlabs/wa-scheduler/internal/core"
	"github.com/go-co-op/gocron/v2"
	"gopkg.in/validator.v2"
)

const (
	// DefaultMaxRetries is the default number of retries for a message
	DefaultMaxRetries = 3
	// DefaultRetryDelay is the default delay between retries
	DefaultRetryDelay = 30 * time.Second
	// DefaultTolerateLateMessage is the default range for tolerating late messages
	DefaultTolerateLateMessage = 1 * time.Minute
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

func (s *GoCronScheduler) ScheduleMessage(ctx context.Context, msg core.Message) error {
	scheduledTime := time.Unix(msg.ScheduledSendingAt, 0)
	now := time.Now()

	// Check if the scheduled time is in the past
	if scheduledTime.Before(now) && scheduledTime.Add(DefaultTolerateLateMessage).Before(now) {
		// If scheduled time is in the past and outside the default range, mark as failed
		reason := "late scheduling message"
		msg.Reason = &reason
		msg.Status = core.MessageStatusFailed

		err := s.Storage.UpdateMessage(ctx, msg)
		if err != nil {
			return fmt.Errorf("failed to update message status to failed: %w", err)
		}

		slog.Error("message scheduled in the past", slog.String("message", msg.String()))
		return nil
	}

	// Check if the scheduled time is within default range (consider as now)
	isNow := scheduledTime.Unix() <= now.Unix()

	var option gocron.OneTimeJobStartAtOption
	if isNow {
		option = gocron.OneTimeJobStartImmediately()
	} else {
		option = gocron.OneTimeJobStartDateTime(scheduledTime)
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
	err := s.Publisher.Publish(ctx, msg)
	if err != nil {
		// if message failed to send more than max retries
		// or if the error is session expired
		// mark the message as failed
		if (msg.RetriedCount >= DefaultMaxRetries) || (err == core.ErrSessionExpired) {
			msg.Status = core.MessageStatusFailed

			reason := fmt.Sprintf("failed to send message after %d retries: %v", DefaultMaxRetries, err)
			if err == core.ErrSessionExpired {
				reason = "session expired"
			}
			msg.Reason = &reason

			slog.Error("failed to send message", slog.String("message", msg.String()), slog.String("err", err.Error()))
		} else {
			// otherwise, retry the message
			msg.ScheduledSendingAt = time.Now().Add(DefaultRetryDelay).Unix()
			s.RetryMessage(ctx, msg)
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

func (s *GoCronScheduler) RetryMessage(ctx context.Context, msg core.Message) error {
	msg.RetriedCount++
	msg.Status = core.MessageStatusScheduled
	err := s.ScheduleMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to retry message: %w", err)
	}

	err = s.Storage.UpdateMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}
