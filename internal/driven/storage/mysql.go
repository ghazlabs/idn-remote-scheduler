package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ghazlabs/idn-remote-scheduler/internal/core"
	"gopkg.in/validator.v2"
)

const (
	tableSchedule = "messages"
)

type MySQLStorageConfig struct {
	DB *sql.DB `validate:"nonnil"`
}

type MySQLStorage struct {
	MySQLStorageConfig
}

func NewMySQLStorage(cfg MySQLStorageConfig) (*MySQLStorage, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return &MySQLStorage{
		MySQLStorageConfig: cfg,
	}, nil
}

func (s *MySQLStorage) GetAllMessages(ctx context.Context, req core.GetAllMessageRequest) ([]core.Message, error) {
	query := fmt.Sprintf("SELECT id, message, scheduled_sending_at, sent_at, retried_count, status, reason, created_at, updated_at FROM %s", tableSchedule)

	var args []interface{}
	var conditions []string

	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []core.Message
	for rows.Next() {
		var msg core.Message
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.ScheduledSendingAt,
			&msg.SentAt,
			&msg.RetriedCount,
			&msg.Status,
			&msg.Reason,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return messages, nil
}

func (s *MySQLStorage) SaveMessage(ctx context.Context, msg core.Message) error {
	query := `
		INSERT INTO messages (
			id, message, scheduled_sending_at, sent_at,
			retried_count, status, reason, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.DB.ExecContext(ctx, query,
		msg.ID,
		msg.Content,
		msg.ScheduledSendingAt,
		msg.SentAt,
		msg.RetriedCount,
		msg.Status,
		msg.Reason,
		msg.CreatedAt,
		msg.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}
