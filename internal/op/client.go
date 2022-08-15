package op

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type Config struct {
	Executable string
}

type Client struct {
	config  Config
	session string
}

func NewClient(cfg Config) *Client {
	return &Client{
		config: cfg,
	}
}

func (c *Client) Exec(args []string) ([]byte, error) {
	cmd := exec.Command(c.config.Executable)
	cmd.Args = append(cmd.Args, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		var tmp *exec.ExitError
		if !errors.As(err, &tmp) {
			return nil, err
		}

		msg, err := ParseExecErrorMessage(stderr.Bytes())
		if err != nil {
			return nil, errors.New(string(stderr.Bytes()))
		}

		return nil, errors.New(msg)
	}
	return stdout.Bytes(), nil
}

var ExecErrorMessageRe = regexp.MustCompile(`^\[ERROR] \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} (.*)(\s.*)*$`)

var ParseExecErrorMessageError = errors.New("parsing error message")

func ParseExecErrorMessage(data []byte) (string, error) {
	match := ExecErrorMessageRe.FindSubmatch(data)
	if len(match) < 2 {
		return "", fmt.Errorf("%w %q", ParseExecErrorMessageError, data)
	}
	return string(match[1]), nil
}

type Secret string

var SecretRe = regexp.MustCompile(`^op:\/(\/[^\/\\]+)+$`)

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var plain string
	if err := unmarshal(&plain); err != nil {
		return err
	}

	if !SecretRe.MatchString(plain) {
		return fmt.Errorf("secret %q must match regex `%s`", plain, SecretRe)
	}

	*s = Secret(plain)
	return nil
}

func (c *Client) Read(secret Secret) ([]byte, error) {
	data, err := c.Exec([]string{
		"read", "-n", string(secret),
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
