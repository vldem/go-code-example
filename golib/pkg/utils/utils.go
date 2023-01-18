//go:generate mockgen -source=./utils.go -destination=./mocks/utils.go -package=mock_utils
package utils

import (
	"math/rand"
	"time"
)

type UtilsInterface interface {
	GetRandIdString(length int) string
}

type utilsImplementation struct{}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func New() UtilsInterface {
	return &utilsImplementation{}
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (u *utilsImplementation) GetRandIdString(length int) string {
	rand.Seed(time.Now().UnixNano())
	return stringWithCharset(length, charset)
}
