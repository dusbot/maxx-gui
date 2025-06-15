package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

func GenerateTimestampUUID(length int) (string, error) {
	if length <= 0 {
		length = 8
	}

	now := time.Now()
	timePart := fmt.Sprintf(
		"%04d%02d%02d-%02d%02d%02d%03d%03d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond()/1e6, now.Nanosecond()/1e3%1e3,
	)
	randomPart, err := generateRandomString(length)
	return timePart + "-" + randomPart, err
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), err
}
