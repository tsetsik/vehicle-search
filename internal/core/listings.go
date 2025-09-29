package core

type (
	Listings interface {
		AddItem(listing []byte) error
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

func (ls *listings) AddItem(listing []byte) error {
	return ls.store.AddItem(listing)
}
