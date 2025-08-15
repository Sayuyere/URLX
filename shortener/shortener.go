package shortener

import (
	"math/rand"
)

type Shortener interface {
	Shorten(url string) (short string)
}

type SimpleShortener struct{}

func NewSimpleShortener() *SimpleShortener {
	return &SimpleShortener{}
}

func (s *SimpleShortener) Shorten(url string) string {
	return generateShortURL()
}

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
