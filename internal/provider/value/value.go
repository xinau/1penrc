package value

import (
	"errors"
	"fmt"

	"github.com/xinau/1penrc/internal/op"
	"github.com/xinau/1penrc/internal/provider"
)

var DefaultConfig = Config{}

type Config struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (cfg *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*cfg = DefaultConfig
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(cfg)); err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
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
	return provider.Variables{
		cfg.Name: cfg.Value,
	}, nil
}
