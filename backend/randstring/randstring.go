package randstring

import (
	"math/rand"
	"time"
)

//const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
const upperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func NewFromCharset(charset string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NewFromUpper(length int) string {
	return NewFromCharset(upperCharset, length)
}

func New(length int) string {
	return NewFromCharset(charset, length)
}
