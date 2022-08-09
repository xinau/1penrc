package value

import (
	"errors"
	"fmt"

	"github.com/xinau/1penrc/internal/op"
	"github.com/xinau/1penrc/internal/provider"
)

type Config struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (cfg *Config) Validate() error {
	if cfg.Name == "" {
		return errors.New("name shouldn't be empty")
	}
	if cfg.Value == "" {
		return errors.New("value can't be empty")
	}
	return nil
}

func GetVariables(_ *op.Client, cfg *Config) (provider.Variables, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return provider.Variables{
		cfg.Name: cfg.Value,
	}, nil
}
