package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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

func (s *MySQLStorage) GetAllMessages(ctx context.Context, message core.Message) ([]core.Message, error) {
	query := fmt.Sprintf("SELECT id, message, scheduled_sending_at, sent_at, retried_count, status, reason, created_at, updated_at FROM %s", tableSchedule)

	var args []interface{}
	var conditions []string

	if message.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, message.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	messages := make([]core.Message, 0)
	for rows.Next() {
		var msg core.Message
		var createdAt, updatedAt string
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.ScheduledSendingAt,
			&msg.SentAt,
			&msg.RetriedCount,
			&msg.Status,
			&msg.Reason,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		t, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at timestamp: %w", err)
		}
		msg.CreatedAt = t.Unix()

		t, err = time.Parse("2006-01-02 15:04:05", updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at timestamp: %w", err)
		}
		msg.UpdatedAt = t.Unix()

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return messages, nil
}

func (s *MySQLStorage) SaveMessage(ctx context.Context, message core.Message) error {
	query := fmt.Sprintf("INSERT INTO %s (message, scheduled_sending_at, sent_at, retried_count, status, reason) VALUES (?, ?, ?, ?, ?, ?)", tableSchedule)

	_, err := s.DB.ExecContext(ctx, query,
		message.Content,
		message.ScheduledSendingAt,
		message.SentAt,
		message.RetriedCount,
		message.Status,
		message.Reason,
	)

	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}

func (s *MySQLStorage) UpdateMessage(ctx context.Context, message core.Message) error {
	var (
		setClauses []string
		args       []interface{}
	)

	if message.Content != "" {
		setClauses = append(setClauses, "message = ?")
		args = append(args, message.Content)
	}
	if message.ScheduledSendingAt != 0 {
		setClauses = append(setClauses, "scheduled_sending_at = ?")
		args = append(args, message.ScheduledSendingAt)
	}
	if message.SentAt != nil && *message.SentAt != 0 {
		setClauses = append(setClauses, "sent_at = ?")
		args = append(args, message.SentAt)
	}
	if message.RetriedCount != 0 {
		setClauses = append(setClauses, "retried_count = ?")
		args = append(args, message.RetriedCount)
	}
	if message.Status != "" {
		setClauses = append(setClauses, "status = ?")
		args = append(args, message.Status)
	}
	if message.Reason != nil {
		setClauses = append(setClauses, "reason = ?")
		args = append(args, message.Reason)
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableSchedule, strings.Join(setClauses, ", "))
	args = append(args, message.ID)

	_, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}
