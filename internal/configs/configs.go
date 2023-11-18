package configs

import (
	"fmt"
	"os"

	"github.com/andrdru/go-template/configs"
	"gopkg.in/yaml.v3"
)

type (
	// Config main config
	Config struct {
		IsDebug  bool             `yaml:"is_debug"`
		Postgres configs.Postgres `yaml:"postgres"`
		HTTP     HTTP             `yaml:"http"`
	}

	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
)

// NewConfig read config from file
func NewConfig(path string) (config Config, err error) {
	var bytes []byte
	bytes, err = os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("could not read file: %w", err)
	}

	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return config, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return config, nil
}
