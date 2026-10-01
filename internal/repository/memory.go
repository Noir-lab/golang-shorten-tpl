package repository

import (
	"golang-shorten-tpl/internal/service"
	"sync"
)

type memoryRepository struct {
	mu   sync.RWMutex
	urls map[string]string
}

var _ service.Repository = (*memoryRepository)(nil)

func NewMemoryRepository() service.Repository {
	return &memoryRepository{
		urls: make(map[string]string),
	}
}

func (r *memoryRepository) Save(id string, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.urls[id] = originalURL
	return nil
}

func (r *memoryRepository) Get(id string) (originalURL string, err bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	originalURL, found := r.urls[id]
	return originalURL, found
}
