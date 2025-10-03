package core

type (
	Listings interface {
		AddItem(items Item) error
	}

	listings struct {
		store Store
	}
)

func NewListings(store Store) Listings {
	return &listings{
		store: store,
	}
}

func (ls *listings) AddItem(item Item) error {
	return ls.store.AddItem(item)
}
