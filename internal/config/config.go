package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/xinau/1penrc/internal/op"
	"github.com/xinau/1penrc/internal/provider/awssts"
	"github.com/xinau/1penrc/internal/provider/secret"
	"github.com/xinau/1penrc/internal/provider/value"
)

func Load(data string) (*Config, error) {
	cfg := &Config{}
	// If the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.
	*cfg = DefaultConfig
	return cfg, yaml.UnmarshalStrict([]byte(data), cfg)
}

func LoadFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return Load(string(data))
}

var (
	DefaultConfig = Config{
		ClientConfig: DefaultClientConfig,
	}

	DefaultClientConfig = op.Config{
		Executable: "op",
	}
)

type Config struct {
	ClientConfig       op.Config            `yaml:"client"`
	EnvironmentConfigs []*EnvironmentConfig `yaml:"environments"`
}

type EnvironmentConfig struct {
	Name    string `yaml:"name"`
	Account string `yaml:"account"`

	AWSSTSConfigs []*awssts.Config `yaml:"aws_sts"`
	SecretConfigs []*secret.Config `yaml:"secrets"`
	ValueConfigs  []*value.Config  `yaml:"values"`
}
