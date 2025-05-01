package core

type Message struct {
	ID                 string   `json:"id"`
	Content            string   `json:"content"`
	RecipientNumbers   []string `json:"recipient_numbers"`
	ScheduledSendingAt int64    `json:"scheduled_sending_at"`
	RetriedCount       int      `json:"retried_count"`
	Status             string   `json:"status"`
	SentAt             *string  `json:"sent_at"`
	Reason             *string  `json:"reason"`
	CreatedAt          int64    `json:"created_at"`
	UpdatedAt          int64    `json:"updated_at"`
}
