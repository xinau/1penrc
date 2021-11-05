package op

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Client struct {
	Executable string
	Account    string
	Session    string
}

func NewClient() *Client {
	return &Client{
		Executable: "op",
	}
}

type ExecError struct {
	Cmd string
	Msg string
	Err error
}

func (e *ExecError) Error() string {
	return e.Msg
}

func (e *ExecError) Is(target error) bool {
	val, ok := target.(*ExecError)
	if !ok {
		return false
	}

	return e.Cmd == val.Cmd && e.Msg == val.Msg
}

func (e *ExecError) Unwrap() error {
	return e.Err
}

func (c *Client) Exec(args []string) ([]byte, error) {
	path, err := exec.LookPath(c.Executable)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(path)

	if c.Session != "" {
		cmd.Args = append(cmd.Args, "--session", c.Session)
	}
	cmd.Args = append(cmd.Args, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		var tmp *exec.ExitError
		if !errors.As(err, &tmp) {
			return nil, &ExecError{
				Cmd: ObfuscateCmd(cmd),
				Msg: err.Error(),
				Err: err,
			}
		}

		msg, err := ParseErrorMsg(stderr.Bytes())
		if err != nil {
			msg = stderr.String()
		}

		return nil, &ExecError{
			Cmd: ObfuscateCmd(cmd),
			Msg: msg,
			Err: err,
		}
	}

	return stdout.Bytes(), nil
}

func Exec(args []string) ([]byte, error) {
	return NewClient().Exec(args)
}

func (c *Client) SignIn() error {
	data, err := c.Exec([]string{"signin", c.Account, "--raw"})
	if err != nil {
		return err
	}

	c.Session, err = ParseSessionToken(data)
	return err
}

func SignIn(account string) (*Client, error) {
	client := NewClient()
	client.Account = account
	return client, client.SignIn()
}

var SessionTokenRe = regexp.MustCompile(`^([0-9A-Za-z_-]{43})\s*$`)

var ParseSessionTokenError = errors.New("parsing session token")

func ParseSessionToken(data []byte) (string, error) {
	match := SessionTokenRe.FindSubmatch(data)
	if len(match) < 2 {
		return "", fmt.Errorf("%w from %q", ParseSessionTokenError, data)
	}
	return string(match[1]), nil
}

var ErrorMsgRe = regexp.MustCompile(`^\[ERROR] \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} (.*)(\s.*)*$`)

var ParsingErrorMsgError = errors.New("parsing error message")

func ParseErrorMsg(data []byte) (string, error) {
	match := ErrorMsgRe.FindSubmatch(data)
	if len(match) < 2 {
		return "", fmt.Errorf("%w %q", ParsingErrorMsgError, data)
	}
	return string(match[1]), nil
}

func ObfuscateCmd(cmd *exec.Cmd) string {
	args := make([]string, len(cmd.Args))
	for i := 0; i < len(cmd.Args); i++ {
		args[i] = cmd.Args[i]
		if cmd.Args[i] == "--session" && i+1 < len(args) {
			i++
			args[i] = "<token>"
		}
	}

	return fmt.Sprintf("%s %s", cmd.Path, strings.Join(args[1:], " "))
}
