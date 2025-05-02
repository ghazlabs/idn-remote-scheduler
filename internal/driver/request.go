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
	Message            string   `json:"message"`
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

func prepareSendMessageRequest(req SendMessageRequest) (SendMessageRequest, error) {
	if len(req.RecipientNumbers) == 0 {
		return req, fmt.Errorf("recipient numbers cannot be empty")
	}

	if req.Message == "" {
		return req, fmt.Errorf("message cannot be empty")
	}

	if req.ScheduledSendingAt == 0 {
		return req, fmt.Errorf("schedule sending at cannot be empty")
	}

	return req, nil
}
