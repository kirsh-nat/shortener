package memoryRepository

import (
	"sync"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewMemoryRepository() services.URLRepository {
	return &MemoryRepository{store: make(map[string]string)}
}

func (r *MemoryRepository) Ping() error {
	return nil
}
