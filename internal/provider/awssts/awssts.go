package awssts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/xinau/1penrc/internal/duration"
	"github.com/xinau/1penrc/internal/op"
	"github.com/xinau/1penrc/internal/provider"
)

var DefaultConfig = Config{
	TTL: duration.Duration(time.Hour),
}

type Config struct {
	TTL             duration.Duration `yaml:"ttl"`
	RoleARN         string            `yaml:"role_arn"`
	AccessKeyID     op.Secret         `yaml:"access_key_id"`
	SecretAccessKey op.Secret         `yaml:"secret_access_key"`
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
	if cfg.RoleARN == "" {
		return errors.New("role_arn can't be empty")
	}
	if cfg.AccessKeyID == "" {
		return errors.New("access_key_id can't be empty")
	}
	if cfg.SecretAccessKey == "" {
		return errors.New("secret_access_key can't be empty")
	}
	if cfg.TTL <= 0 {
		return errors.New("ttl must be greater than zero")
	}
	return nil
}

func GetVariables(client *op.Client, cfg *Config) (provider.Variables, error) {
	svc := sts.New(sts.Options{
		Credentials: NewOnePasswordProvider(client, cfg),
		Region:      "us-east-1",
	})

	creds, err := stscreds.NewAssumeRoleProvider(
		svc, cfg.RoleARN,
		func(opts *stscreds.AssumeRoleOptions) {
			if cfg.TTL > 0 {
				opts.Duration = time.Duration(cfg.TTL)
			}
		},
	).Retrieve(context.TODO())

	if err != nil {
		return nil, err
	}

	return provider.Variables{
		"AWS_ACCESS_KEY_ID":     creds.AccessKeyID,
		"AWS_SECRET_ACCESS_KEY": creds.SecretAccessKey,
		"AWS_SESSION_TOKEN":     creds.SessionToken,
	}, nil
}

type OnePasswordProvider struct {
	config *Config
	client *op.Client
}

func NewOnePasswordProvider(client *op.Client, cfg *Config) *OnePasswordProvider {
	return &OnePasswordProvider{
		config: cfg,
		client: client,
	}
}

func (p *OnePasswordProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	key, err := p.client.Read(p.config.AccessKeyID)
	if err != nil {
		return aws.Credentials{}, err
	}

	secret, err := p.client.Read(p.config.SecretAccessKey)
	if err != nil {
		return aws.Credentials{}, err
	}

	return aws.Credentials{
		AccessKeyID:     string(key),
		SecretAccessKey: string(secret),
		Source:          "1Password",
	}, nil
}
