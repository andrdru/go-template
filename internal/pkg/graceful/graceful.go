package graceful

import (
	"context"
	"fmt"
	"log/slog"
)

type (
	Closer func(ctx context.Context) (description string, err error)
)

// Stop close app resources with LIFO order
func Stop(ctx context.Context, logger *slog.Logger, closers []Closer) {
	if len(closers) == 0 {
		return
	}

	for i := len(closers) - 1; i >= 0; i-- {
		description, errClose := closers[i](ctx)
		if errClose != nil {
			logger.Error(fmt.Sprintf("close: %s", description), slog.Any("error", errClose))
		} else {
			logger.Info(fmt.Sprintf("closed: %s", description))
		}
	}
}
