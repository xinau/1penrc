package secret

import (
	"errors"
	"fmt"

	"github.com/xinau/1penrc/internal/op"
	"github.com/xinau/1penrc/internal/provider"
)

type Config struct {
	Name   string    `yaml:"name"`
	Secret op.Secret `yaml:"secret"`
}

func (cfg *Config) Validate() error {
	if cfg.Name == "" {
		return errors.New("name shouldn't be empty")
	}
	if cfg.Secret == "" {
		return errors.New("secret can't be empty")
	}
	return nil
}

func GetVariables(client *op.Client, cfg *Config) (provider.Variables, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	val, err := client.Read(cfg.Secret)
	if err != nil {
		return nil, err
	}

	return provider.Variables{
		cfg.Name: string(val),
	}, nil
}
