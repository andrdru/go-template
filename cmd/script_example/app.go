package script_example

import (
	"log/slog"
)

func Run(logger *slog.Logger) (code int) {
	logger.Info("hello, I'm script example")

	return 0
}
