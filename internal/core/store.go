package core

//go:generate mockgen -source=./store.go -destination=./mocks/store_mock.go -package=mocks

import (
	"encoding/base64"
	"encoding/json"

	"github.com/tsetsik/vehicle-search/internal/cache"
)

type (
	Store interface {
		AddItem(item any) error
		GetItem(hash string) any
	}

	store struct {
		cache cache.CacheStore
	}
)

func NewStore(cache cache.CacheStore) Store {
	return &store{
		cache: cache,
	}
}

func (s *store) AddItem(item any) error {
	s.cache.Put(s.Hash(item), item)
	return nil
}

func (s *store) GetItem(hash string) any {
	item, exists := s.cache.Get(hash)
	if !exists {
		return nil
	}

	return item
}

func (s *store) Hash(item any) string {
	// Implement a proper hashing function here
	jstr, _ := json.Marshal(item)
	return base64.StdEncoding.EncodeToString(jstr)
}
