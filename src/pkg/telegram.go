package pkg

import (
	"crazyfarmbackend/src/domain/dto"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrAuthDateMissing  = errors.New("auth_date is missing")
	ErrSignMissing      = errors.New("sign is missing")
	ErrSignInvalid      = errors.New("sign is invalid")
	ErrUnexpectedFormat = errors.New("init data has unexpected format")
	ErrExpired          = errors.New("init data is expired")
)

var (
	_stringProps = map[string]bool{
		"start_param": true,
	}
)

func ParseTelegramData(initData string) (dto.InitData, error) {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return dto.InitData{}, ErrUnexpectedFormat
	}
	pairs := make([]string, 0, len(q))
	for k, v := range q {
		val := v[0]
		valFormat := "%q:%q"
		if isString := _stringProps[k]; !isString && json.Valid([]byte(val)) {
			valFormat = "%q:%s"
		}

		pairs = append(pairs, fmt.Sprintf(valFormat, k, val))
	}
	var d dto.InitData
	jStr := fmt.Sprintf("{%s}", strings.Join(pairs, ","))
	if err := json.Unmarshal([]byte(jStr), &d); err != nil {
		return dto.InitData{}, ErrUnexpectedFormat
	}
	return d, nil
}

func ValidateTelegramData(initData, token string, expIn time.Duration) error {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return ErrUnexpectedFormat
	}

	var (
		authDate time.Time
		hash     string
		pairs    = make([]string, 0, len(q))
	)
	for k, v := range q {
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "auth_date" {
			if i, err := strconv.Atoi(v[0]); err == nil {
				authDate = time.Unix(int64(i), 0)
			}
		}
		pairs = append(pairs, k+"="+v[0])
	}
	if hash == "" {
		return ErrSignMissing
	}
	if expIn > 0 {
		if authDate.IsZero() {
			return ErrAuthDateMissing
		}
		if authDate.Add(expIn).Before(time.Now()) && os.Getenv("APP") == "release" {
			return ErrExpired
		}
	}
	sort.Strings(pairs)
	if sign(strings.Join(pairs, "\n"), token) != hash && os.Getenv("APP") == "release" {
		return ErrSignInvalid
	}
	return nil
}

func SignTelegramData(payload map[string]string, key string, authDate time.Time) string {
	pairs := make([]string, 0, len(payload)+1)
	for k, v := range payload {
		if k == "hash" || k == "auth_date" {
			continue
		}
		pairs = append(pairs, k+"="+v)
	}

	pairs = append(pairs, "auth_date="+strconv.FormatInt(authDate.Unix(), 10))
	sort.Strings(pairs)
	return sign(strings.Join(pairs, "\n"), key)
}

func SignQueryString(qs, key string, authDate time.Time) (string, error) {
	qp, err := url.ParseQuery(qs)
	if err != nil {
		return "", err
	}
	m := make(map[string]string, len(qp))
	for k, v := range qp {
		m[k] = v[0]
	}
	return SignTelegramData(m, key, authDate), nil
}

func sign(payload, key string) string {
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(key))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(payload))

	return hex.EncodeToString(impHmac.Sum(nil))
}
