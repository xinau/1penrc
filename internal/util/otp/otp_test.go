package otp

import (
	"errors"
	"testing"
	"time"
)

func TestParseSecretFromURL(t *testing.T) {
	table := []struct {
		name    string
		data    string
		want    string
		wantErr error
	}{{
		"otpauth uri",
		`otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP`,
		`JBSWY3DPEHPK3PXP`,
		nil,
	}}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			seed, err := ParseSecretFromURL(test.data)
			Assertf(t, errors.Is(err, test.wantErr), "got %s, expected %s", err, test.wantErr)
			Assertf(t, seed == test.want, "got %q, expected %q", seed, test.want)
		})
	}
}

func TestGenerateTOTPCode(t *testing.T) {
	table := []struct {
		name    string
		secret  string
		want    string
		wantErr error
	}{{
		"simple totp",
		`JBSWY3DPEHPK3PXP`,
		`282760`,
		nil,
	}}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			code, err := GenerateTOTPCode(test.secret, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
			Assertf(t, errors.Is(err, test.wantErr), "got %s, expected %s", err, test.wantErr)
			Assertf(t, code == test.want, "got %q, expected %q", code, test.want)
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
