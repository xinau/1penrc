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
	MultiFactorAuth MultiFactorAuth   `yaml:",inline"`
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

func (cfg *Config) HasMultiFactorAuth() bool {
	return cfg.MultiFactorAuth.Token != ""
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

type MultiFactorAuth struct {
	Token        op.Secret `mfa_token`
	SerialNumber string    `mfa_serial_number`
}

func GetVariables(client *op.Client, cfg *Config) (provider.Variables, error) {
	var secret, serialNumber string
	if cfg.HasMultiFactorAuth() {
		data, err := client.Read(cfg.MultiFactorAuth.Token)
		if err != nil {
			return nil, err
		}

		secret, err = ParseOTPSecretFromURL(string(data))
		if err != nil {
			return nil, err
		}

		serialNumber = cfg.MultiFactorAuth.SerialNumber
	}

	creds := stscreds.NewAssumeRoleProvider(
		sts.New(sts.Options{
			Credentials: NewOnePasswordProvider(client, cfg),
			Region:      "us-east-1",
		}),
		cfg.RoleARN,
		func(opts *stscreds.AssumeRoleOptions) {
			if cfg.TTL > 0 {
				opts.Duration = time.Duration(cfg.TTL)
			}
		},
		func(opts *stscreds.AssumeRoleOptions) {
			if cfg.HasMultiFactorAuth() {
				opts.SerialNumber = &serialNumber
			}
		},
		func(opts *stscreds.AssumeRoleOptions) {
			if cfg.HasMultiFactorAuth() {
				opts.TokenProvider = TOTPTokenProvider(secret)
			}
		},
	)

	awscreds, err := creds.Retrieve(context.TODO())
	if err != nil {
		return nil, err
	}

	return provider.Variables{
		"AWS_ACCESS_KEY_ID":     awscreds.AccessKeyID,
		"AWS_SECRET_ACCESS_KEY": awscreds.SecretAccessKey,
		"AWS_SESSION_TOKEN":     awscreds.SessionToken,
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

func TOTPTokenProvider(secret string) func() (string, error) {
	return func() (string, error) {
		return GenerateTOTPCode(secret, time.Now())
	}
}
