package token

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const separator = ":"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Token struct {
	User      string
	Expiry    time.Time
	Signature []byte
}

func (t *Token) Validate(secret []byte) error {
	if len(secret) < 32 {
		panic("invalid secret")
	}
	if t.Expiry.Before(time.Now()) {
		return ErrTokenExpired
	}
	return nil
}

func Generate(secret []byte, duration time.Duration, user string) (string, error) {
	expiry := time.Now().Add(duration).UnixMilli()
	head := fmt.Sprintf("%s%s%d", user, separator, expiry)
	hasher := hmac.New(sha512.New, secret)
	_, err := io.WriteString(hasher, head)
	if err != nil {
		return "", err
	}
	signature := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("%s%s%s", head, separator, signature), nil
}

func Parse(input string) (*Token, error) {
	parts := strings.Split(input, separator)
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	expiryUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	sigBytes, err := base64.URLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Token{
		User:      parts[0],
		Expiry:    time.UnixMilli(expiryUnix),
		Signature: sigBytes,
	}, nil
}
