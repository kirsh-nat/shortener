package memoryrepository

import (
	"sync"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

type UserDataURL struct {
	UserID      string
	OriginalURL string
	Deleted     bool
}

type MemoryRepository struct {
	mu    sync.RWMutex
	store map[string]UserDataURL
}

func NewMemoryRepository() services.URLRepository {
	return &MemoryRepository{store: make(map[string]UserDataURL)}
}

func (r *MemoryRepository) Ping() error {
	return nil
}
