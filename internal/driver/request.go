package driver

import (
	"fmt"

	"gopkg.in/validator.v2"
)

type GetAllMessageRequest struct {
	Status *string
}

type SendMessageRequest struct {
	RecipientNumbers   []string `json:"recipient_numbers"`
	Content            string   `json:"content"`
	ScheduledSendingAt int64    `json:"scheduled_sending_at"`
}

func (r SendMessageRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return fmt.Errorf("invalid vacancy: %w", err)
	}
	return nil
}

type RetryMessageRequest struct {
	ScheduledSendingAt int64 `json:"scheduled_sending_at"`
}

func (r RetryMessageRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return fmt.Errorf("invalid vacancy: %w", err)
	}
	return nil
}
