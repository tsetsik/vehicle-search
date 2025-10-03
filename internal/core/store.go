package core

//go:generate mockgen -source=./store.go -destination=./mocks/store_mock.go -package=mocks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/tsetsik/vehicle-search/internal/cache"
)

type (
	Store interface {
		AddItem(item any) error
		GetItem(hash string) any
		GetAllItems() []any
	}

	store struct {
		cache cache.CacheStore
	}
)

func NewStore[A any](cache cache.CacheStore) Store {
	return &store{
		cache: cache,
	}
}

func (s *store) AddItem(item any) error {
	hash, err := s.Hash(item)
	if err != nil {
		return fmt.Errorf("failed to hash item: %w", err)
	}

	s.cache.Put(hash, item)
	return nil
}

func (s *store) GetItem(hash string) any {
	item, exists := s.cache.Get(hash)
	if !exists {
		return nil
	}

	return item
}

func (s *store) GetAllItems() []any {
	items := s.cache.List()
	result := make([]any, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result
}

func (s *store) Hash(item any) (string, error) {
	// Implement a proper hashing function here
	jstr, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("failed to marshal item: %w", err)
	}

	return base64.StdEncoding.EncodeToString(jstr), nil
}
