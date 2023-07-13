package script_example

import (
	"github.com/rs/zerolog"
)

func Run(logger zerolog.Logger) (code int) {
	logger.Info().Msgf("hello, I'm script example")

	return 0
}
