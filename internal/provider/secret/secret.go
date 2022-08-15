package secret

import (
	"errors"
	"fmt"

	"github.com/xinau/1penrc/internal/op"
	"github.com/xinau/1penrc/internal/provider"
)

var DefaultConfig = Config{}

type Config struct {
	Name   string    `yaml:"name"`
	Secret op.Secret `yaml:"secret"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
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
	val, err := client.Read(cfg.Secret)
	if err != nil {
		return nil, err
	}

	return provider.Variables{
		cfg.Name: string(val),
	}, nil
}
