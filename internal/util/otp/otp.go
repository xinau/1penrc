package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

func ParseSecretFromURL(data string) (string, error) {
	str := strings.TrimSpace(data)
	uri, err := url.ParseRequestURI(str)
	if err != nil {
		return "", err
	}

	vals := uri.Query()
	if !vals.Has("secret") {
		return "", errors.New("missing secret in otpauth uri")
	}

	secret := vals.Get("secret")
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	return strings.ToUpper(secret), nil
}

func GenerateTOTPCode(secret string, date time.Time) (string, error) {
	counter := uint64(math.Floor(float64(date.Unix()) / float64(30)))

	data, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", errors.New("secret isn't base32 encoded")
	}

	buf := make([]byte, 8)
	mac := hmac.New(sha1.New, data)
	binary.BigEndian.PutUint64(buf, counter)
	mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	mod := int32(value % int64(math.Pow10(6)))
	return fmt.Sprintf("%06d", mod), nil
}
