package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/andrdru/go-template/cmd/app"
	"github.com/andrdru/go-template/cmd/script_example"
)

type (
	flags struct {
		isHelp     *bool
		configPath *string
		script     *string
	}
)

const (
	// serviceName name of service
	// redefine here or with ldflags
	// go build -ldflags="-X 'main.serviceName=my_service'"
	serviceName = "service"

	scriptExample = "example"
)

func main() {
	f := initFlags()
	if *f.isHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	logger := initLogger()

	var (
		code int
	)

	switch *f.script {
	case "":
		code = app.Run(logger, *f.configPath)
	default:
		logger.Error().Msgf("unknown script: %s", *f.script)
		os.Exit(1)

	case scriptExample:
		code = script_example.Run(logger)
	}

	os.Exit(code)
}

func initFlags() (fv flags) {
	fv.isHelp = flag.Bool("help", false, "Print help and exit")
	fv.configPath = flag.String("config", "config.yaml", "path to config.yml")
	fv.script = flag.String("script", "", "Run in script mode. One of: test-script")

	flag.Parse()
	return fv
}

func initLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	return zerolog.New(os.Stdout).
		With().Timestamp().
		Str("service", serviceName).
		Logger()
}
