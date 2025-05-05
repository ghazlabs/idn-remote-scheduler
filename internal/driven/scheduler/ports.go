package scheduler

import (
	"context"

	"github.com/ghazlabs/idn-remote-scheduler/internal/core"
)

type Publisher interface {
	Publish(ctx context.Context, msg core.Message) error
}
