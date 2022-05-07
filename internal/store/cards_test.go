package store

import (
	"reflect"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestStore_GetCard(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get single card#", func(t *testing.T) {
		user := "user1"

		card := common.Card{
			Name:   "card 1",
			Number: "1111 2222 3333 4444",
		}

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreCard(user, card)
		assert.NoError(t, err)

		cardRet, err := store.GetCard(user, id)
		assert.NoError(t, err)

		assert.True(t, reflect.DeepEqual(card, cardRet))
	})
}

func TestStore_ListCards(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Get multiple cards", func(t *testing.T) {
		user := "user1"
		name1 := "card 1"
		name2 := "card 2"
		number1 := "1111 1111 1111 1111"
		number2 := "2222 2222 2222 2222"

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id1, err := store.StoreCard(user, common.Card{
			Name:   name1,
			Number: number1,
		})
		assert.NoError(t, err)

		id2, err := store.StoreCard(user, common.Card{
			Name:   name2,
			Number: number2,
		})
		assert.NoError(t, err)

		cards, err := store.ListCards(user)
		assert.NoError(t, err)

		wantCards := make(common.Cards)
		wantCards[id1] = common.Card{
			Name: name1,
		}
		wantCards[id2] = common.Card{
			Name: name2,
		}
		assert.True(t, reflect.DeepEqual(cards, wantCards))
	})
}

func TestStore_DeleteCard(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Delete card record", func(t *testing.T) {
		user := "user1"

		card := common.Card{
			Name:   "the card",
			Number: "1111 2222 3333 4444",
		}

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreCard(user, card)
		assert.NoError(t, err)

		err = store.DeleteCard(user, id)
		assert.NoError(t, err)

		// same attempt should result in error
		err = store.DeleteCard(user, id)
		assert.Error(t, err)

		// attempt to delete non-existing card# recourd
		// should result in error
		err = store.DeleteCard(user, 999999)
		assert.Error(t, err)
	})
}

func TestStore_UpdateCard(t *testing.T) {
	store := dropCreateStore(t)
	t.Run("Update card record", func(t *testing.T) {
		user := "user1"
		name := "card"
		number1 := "1111 1111 1111 1111"
		number2 := "2222 2222 2222 2222"

		_, err := store.AddUser(common.User{
			Name: user,
		})
		assert.NoError(t, err)

		id, err := store.StoreCard(user, common.Card{
			Name:   name,
			Number: number1,
		})
		assert.NoError(t, err)

		err = store.UpdateCard(user, id, common.Card{
			Name:   name,
			Number: number2,
		})
		assert.NoError(t, err)

		card, err := store.GetCard(user, id)
		assert.NoError(t, err)
		assert.Equal(t, card.Number, number2, "Updated card number should change")
	})
}
