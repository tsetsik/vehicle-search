package cache

import (
	"container/list"
	"sync"
)

type (
	CacheStore interface {
		Get(key string) (any, bool)
		Put(key string, value any)
		List() []any
	}

	cacheStore struct {
		capacity int
		entries  map[string]*list.Element
		list     *list.List
		mu       sync.RWMutex
	}
)

func NewLRUCacheStore(capacity int) CacheStore {
	return &cacheStore{
		entries:  make(map[string]*list.Element),
		capacity: capacity,
		list:     list.New(),
	}
}

func (s *cacheStore) Get(key string) (any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.entries[key]
	if !exists {
		return nil, false
	}

	s.list.MoveToFront(e)

	return e.Value, true
}

func (s *cacheStore) List() []any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.list.Len() == 0 {
		return nil
	}

	element := s.list.Front()
	results := []any{element.Value}
	for e := element.Next(); e != nil; e = e.Next() {
		val := e.Value
		results = append(results, val)
		element = e
	}

	return results
}

func (s *cacheStore) Put(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.entries) >= s.capacity {
		back := s.list.Back()
		if back != nil {
			s.list.Remove(back)
			delete(s.entries, key)
		}
	}

	if e, exists := s.entries[key]; exists {
		s.list.MoveToFront(e)
		e.Value = value
		return
	}

	element := s.list.PushFront(value)
	s.entries[key] = element
}
