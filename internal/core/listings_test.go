package core_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tsetsik/vehicle-search/internal/core"
	"github.com/tsetsik/vehicle-search/internal/core/mocks"
)

func Test_AddItem(t *testing.T) {
	t.Parallel()

	t.Run("success on add item to listings", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mocks.NewMockStore(ctrl)
		item := core.Item{
			Make:        "Toyota",
			Model:       "Camry",
			Year:        2020,
			FuelType:    "diesel",
			Description: "A reliable sedan",
		}
		store.EXPECT().AddItem(item).Return(nil).Times(1)

		listings := core.NewListings(store)

		err := listings.AddItem(item)
		require.NoError(t, err)
	})

	t.Run("fail on add item to listings", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mocks.NewMockStore(ctrl)
		item := core.Item{
			Make:        "Toyota",
			Model:       "Camry",
			Year:        2020,
			FuelType:    "diesel",
			Description: "A reliable sedan",
		}
		itemErr := fmt.Errorf("failed to add item")
		store.EXPECT().AddItem(item).Return(itemErr).Times(1)

		listings := core.NewListings(store)

		err := listings.AddItem(item)
		require.Error(t, err, itemErr)
	})
}
