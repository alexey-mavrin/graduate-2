package store

import (
	"database/sql"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	// sqlite sql package
	_ "github.com/mattn/go-sqlite3"
)

// StoreCardID stores card with the ID specified
func (s *Store) StoreCardID(id int64, user string, card common.Card) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	_, err := secretStore.db.Exec(`INSERT INTO cards
		(id, user_id, name, holder, number, expm, expy, cvc, meta)
		VALUES(?, (SELECT id from users where user=?), ?, ?, ?, ?, ?, ?, ?)`,
		id,
		user,
		card.Name,
		card.Holder,
		card.Number,
		card.ExpMonth,
		card.ExpYear,
		card.CVC,
		card.Meta,
	)
	if err != nil {
		return err
	}

	return nil
}

// StoreCard stores Card data for given user
func (s *Store) StoreCard(user string, card common.Card) (int64, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`INSERT INTO cards
		(user_id, name, holder, number, expm, expy, cvc, meta)
		VALUES((SELECT id from users where user=?), ?, ?, ?, ?, ?, ?, ?)`,
		user,
		card.Name,
		card.Holder,
		card.Number,
		card.ExpMonth,
		card.ExpYear,
		card.CVC,
		card.Meta,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateCard updates an card record
func (s *Store) UpdateCard(user string, id int64, card common.Card) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`UPDATE cards
		SET name = ?, holder = ?, number = ?, expm = ?,
		expy = ?, cvc = ?, meta = ?
		WHERE id in
		( SELECT cards.id FROM cards
			JOIN users ON cards.user_id = users.id
			WHERE users.user = ? AND cards.id = ?)`,
		card.Name,
		card.Holder,
		card.Number,
		card.ExpMonth,
		card.ExpYear,
		card.CVC,
		card.Meta,
		user,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrNotFound
	}

	return nil
}

// ListCards returns list of stored cards for the given user
// with all but password fileds filled
func (s *Store) ListCards(user string) (common.Cards, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	cards := make(common.Cards)
	rows, err := secretStore.db.Query(
		`SELECT cards.id, cards.name
			FROM cards JOIN users ON cards.user_id = users.id
			WHERE users.user = ?`,
		user,
	)
	if err != nil {
		return cards, err
	}

	for rows.Next() {
		var id int64
		var card common.Card
		err = rows.Scan(&id, &card.Name)
		if err != nil {
			return cards, err
		}
		cards[id] = card
	}
	return cards, nil
}

// GetCard returns stored card record including the password
func (s *Store) GetCard(user string, id int64) (common.Card, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	var card common.Card

	row := secretStore.db.QueryRow(
		`SELECT cards.name, cards.holder, cards.number, cards.expm,
			cards.expy, cards.cvc, cards.meta
			FROM cards JOIN users ON cards.user_id = users.id
			WHERE users.user = ? AND cards.id = ?`,
		user, id,
	)

	err := row.Scan(&card.Name,
		&card.Holder,
		&card.Number,
		&card.ExpMonth,
		&card.ExpYear,
		&card.CVC,
		&card.Meta,
	)
	if err == sql.ErrNoRows {
		return card, ErrNotFound
	}
	if err != nil {
		return card, err
	}
	return card, nil
}

// DeleteCard deletes the specified card record
func (s *Store) DeleteCard(user string, id int64) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(
		`DELETE FROM cards
		 	WHERE id in
			( SELECT cards.id FROM cards
				JOIN users ON cards.user_id = users.id
				WHERE users.user = ? AND cards.id = ?
			)`,
		user, id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil
	}
	if rows != 1 {
		return ErrNotFound
	}
	return nil
}
