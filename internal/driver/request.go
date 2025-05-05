package driver

type GetAllMessageRequest struct {
	Status *string
}

type SendMessageRequest struct {
	RecipientNumbers   []string `json:"recipient_numbers"`
	Content            string   `json:"content"`
	ScheduledSendingAt int64    `json:"scheduled_sending_at"`
}

type RetryMessageRequest struct {
	ScheduledSendingAt int64 `json:"scheduled_sending_at"`
}
