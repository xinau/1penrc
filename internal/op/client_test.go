package op

import (
	"errors"
	"fmt"
	"os/exec"
	"testing"
)

func TestCommand_Exec(t *testing.T) {
	tests := []struct {
		name string
		cmd  *Client
		args []string
		want string
		err  error
	}{{
		"with account set",
		&Client{"./testdata/op-echo.sh", "example", ""},
		[]string{"--version"},
		"--account example --version\n",
		nil,
	}, {
		"with session set",
		&Client{"./testdata/op-echo.sh", "", "<token>"},
		[]string{"--version"},
		"--session <token> --version\n",
		nil,
	}, {
		"execution error",
		&Client{Executable: "./testdata/op-error.sh"},
		[]string{ErrorMsgSiginUnauthorzied},
		"",
		&ExecError{
			Cmd: fmt.Sprintf("./testdata/op-error.sh %s", ErrorMsgSiginUnauthorzied),
			Msg: ErrorMsgSiginUnauthorzied,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.cmd.Exec(test.args)
			Assertf(t, string(got) == test.want, "got %q, expected %q", got, test.want)
			Assertf(t, errors.Is(err, test.err), "got %q, expected %q", got, test.err)
		})
	}
}

var (
	ErrorMsgSiginUnauthorzied = "(401) Unauthorized: You aren't authorized to perform this action."
	ErrorMsgItemUnknown       = "\"unknown\" doesn't seem to be an item. Specify the item with its UUID, name, or domain."
)

func OpErrorMsg(msg string) string {
	return fmt.Sprintf("[ERROR] 1970/01/01 00:00:00 %s\n", msg)
}

func TestParseErrorMsg(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
		err     error
	}{{
		"op signin with empty password",
		OpErrorMsg(ErrorMsgSiginUnauthorzied),
		ErrorMsgSiginUnauthorzied,
		nil,
	}, {
		"op get item with unknown item",
		OpErrorMsg(ErrorMsgItemUnknown),
		ErrorMsgItemUnknown,
		nil,
	}, {
		"invalid error message",
		"[INVALID] an invalid error message\n",
		"",
		ParsingErrorMsgError,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseErrorMsg([]byte(test.message))
			Assertf(t, got == test.want, "got %q, expected %q", got, test.want)
			Assertf(t, errors.Is(err, test.err), "got %q, expected %q", got, test.err)
		})
	}
}

func TestObfuscateCmd(t *testing.T) {
	tests := []struct {
		name string
		cmd  *exec.Cmd
		want string
	}{{
		"cmd without session arg",
		&exec.Cmd{
			Path: "./example",
			Args: []string{"--account", "example", "account"},
		},
		"./example --account example account",
	}, {
		"cmd with session arg",
		&exec.Cmd{
			Path: "./example",
			Args: []string{"--account", "example", "--session", "1234556", "account"},
		},
		"./example --account example --session <token> account",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ObfuscateCmd(test.cmd)
			Assertf(t, got == test.want, "got %q, expected %q", got, test.want)
		})
	}
}

// Assertf errors if the test clause fails with format and args
func Assertf(t *testing.T, clause bool, format string, args ...interface{}) {
	t.Helper()
	if !clause {
		t.Errorf(format, args...)
	}
}
