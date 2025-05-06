package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ghazlabs/wa-scheduler/internal/core"
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

func (s *MySQLStorage) GetAllMessages(ctx context.Context, input core.GetAllMessagesInput) ([]core.Message, error) {
	query := fmt.Sprintf(`SELECT
		id,
		content,
		recipient_numbers,
		scheduled_sending_at,
		sent_at,
		retried_count,
		status,
		reason,
		created_at,
		updated_at
	FROM %s`, tableSchedule)

	var args []interface{}
	var conditions []string

	if input.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, input.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	messages := make([]core.Message, 0)
	for rows.Next() {
		var msg core.Message
		var recipientNumbers, createdAt, updatedAt string
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&recipientNumbers,
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

		msg.RecipientNumbers = strings.Split(recipientNumbers, ",")

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
func (s *MySQLStorage) GetMessage(ctx context.Context, id string) (*core.Message, error) {
	query := fmt.Sprintf(`SELECT
		id,
		content,
		recipient_numbers,
		scheduled_sending_at,
		sent_at,
		retried_count,
		status,
		reason,
		created_at,
		updated_at
	FROM %s WHERE id = ?`, tableSchedule)

	row := s.DB.QueryRowContext(ctx, query, id)

	var msg core.Message
	var recipientNumbers, createdAt, updatedAt string
	err := row.Scan(
		&msg.ID,
		&msg.Content,
		&recipientNumbers,
		&msg.ScheduledSendingAt,
		&msg.SentAt,
		&msg.RetriedCount,
		&msg.Status,
		&msg.Reason,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	msg.RecipientNumbers = strings.Split(recipientNumbers, ",")

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

	return &msg, nil
}

func (s *MySQLStorage) SaveMessage(ctx context.Context, message core.Message) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, content, scheduled_sending_at, recipient_numbers, status)
		VALUES (?, ?, ?, ?, ?)
	`, tableSchedule)

	_, err := s.DB.ExecContext(ctx, query,
		message.ID,
		message.Content,
		message.ScheduledSendingAt,
		strings.Join(message.RecipientNumbers, ","),
		message.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}

func (s *MySQLStorage) UpdateMessage(ctx context.Context, message core.Message) error {
	query := fmt.Sprintf(`
		UPDATE %s
			SET scheduled_sending_at = ?,
				sent_at = ?,
				retried_count = ?,
				status = ?,
				reason = ?
		WHERE id = ?
	`, tableSchedule)

	_, err := s.DB.ExecContext(ctx, query,
		message.ScheduledSendingAt,
		message.SentAt,
		message.RetriedCount,
		message.Status,
		message.Reason,
		message.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}
