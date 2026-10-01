package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
)

const idAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const idLength = 8

type urlService struct {
	repo Repository
}

func isValidURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

func NewURLService(repo Repository) Service {
	return &urlService{repo: repo}
}

func (s *urlService) Shorten(originalURL string) (shortID string, err error) {
	if !isValidURL(originalURL) {
		return "", ErrInvalidURL
	}

	id := generateID()
	err = s.repo.Save(id, originalURL)
	if err != nil {
		return "", fmt.Errorf("save url: %w", err)
	}
	return id, nil
}

func (s *urlService) GetOriginal(id string) (originalURL string, found bool) {
	if id == "" {
		return "", false
	}
	return s.repo.Get(id)
}

func generateID() string {
	result := make([]byte, idLength)
	for i := range result {
		n, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(idAlphabet))),
		)
		if err != nil {
			panic(err)
		}
		result[i] = idAlphabet[n.Int64()]
	}
	return string(result)
}

var ErrInvalidURL = errors.New("invalid URL")
