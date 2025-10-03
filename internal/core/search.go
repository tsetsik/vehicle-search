package core

import (
	"strings"

	"github.com/samber/lo"
)

type (
	SearchEngine interface {
		Search(filter ItemsFilter) ([]*Item, error)
	}

	ItemsFilter struct {
		Q string
	}

	searchEngine struct {
		store Store
	}
)

func NewSearchEngine(store Store) SearchEngine {
	return &searchEngine{
		store: store,
	}
}

func (se *searchEngine) Search(filter ItemsFilter) ([]*Item, error) {
	// Implement search logic here
	items := se.store.GetAllItems()
	i := lo.Map(items, func(item any, _ int) *Item {
		i, ok := item.(Item)
		if !ok {
			return nil
		}
		return &i
	})

	results := i

	if filter.Q != "" {
		results = lo.Filter(lo.Compact(i), func(item *Item, _ int) bool {
			return strings.Contains(strings.ToLower(item.Description), strings.ToLower(filter.Q)) ||
				strings.Contains(strings.ToLower(item.Model), strings.ToLower(filter.Q))
		})
	}

	return results, nil
}
