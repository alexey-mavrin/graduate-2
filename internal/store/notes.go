package store

import (
	"database/sql"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	// sqlite sql package
	_ "github.com/mattn/go-sqlite3"
)

// StoreNoteID stores note with the ID specified
func (s *Store) StoreNoteID(id int64, user string, note common.Note) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	_, err := secretStore.db.Exec(`INSERT INTO notes
		(id, user_id, name, text, meta)
		VALUES(?, (SELECT id from users where user=?), ?, ?, ?)`,
		id,
		user,
		note.Name,
		note.Text,
		note.Meta,
	)
	if err != nil {
		return err
	}

	return nil
}

// StoreNote stores Note data for given user
func (s *Store) StoreNote(user string, note common.Note) (int64, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`INSERT INTO notes
		(user_id, name, text, meta)
		VALUES((SELECT id from users where user=?), ?, ?, ?)`,
		user,
		note.Name,
		note.Text,
		note.Meta,
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

// UpdateNote updates an note record
func (s *Store) UpdateNote(user string, id int64, note common.Note) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`UPDATE notes
		SET name = ?, text = ?, meta = ?
		WHERE id in
		( SELECT notes.id FROM notes
			JOIN users ON notes.user_id = users.id
			WHERE users.user = ? AND notes.id = ?)`,
		note.Name,
		note.Text,
		note.Meta,
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

// ListNotes returns list of stored notes for the given user
// with all but password fileds filled
func (s *Store) ListNotes(user string) (common.Notes, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	notes := make(common.Notes)
	rows, err := secretStore.db.Query(
		`SELECT notes.id, notes.name, notes.text, notes.meta
			FROM notes JOIN users ON notes.user_id = users.id
			WHERE users.user = ?`,
		user,
	)
	if err != nil {
		return notes, err
	}

	for rows.Next() {
		var id int64
		var note common.Note
		err = rows.Scan(&id, &note.Name, &note.Text, &note.Meta)
		if err != nil {
			return notes, err
		}
		notes[id] = note
	}
	return notes, nil
}

// GetNote returns stored note record including the password
func (s *Store) GetNote(user string, id int64) (common.Note, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	var note common.Note

	row := secretStore.db.QueryRow(
		`SELECT notes.name, notes.text, notes.meta
			FROM notes JOIN users ON notes.user_id = users.id
			WHERE users.user = ? AND notes.id = ?`,
		user, id,
	)

	err := row.Scan(&note.Name,
		&note.Text,
		&note.Meta,
	)
	if err == sql.ErrNoRows {
		return note, ErrNotFound
	}
	if err != nil {
		return note, err
	}
	return note, nil
}

// DeleteNote deletes the specified note record
func (s *Store) DeleteNote(user string, id int64) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(
		`DELETE FROM notes
		 	WHERE id in
			( SELECT notes.id FROM notes
				JOIN users ON notes.user_id = users.id
				WHERE users.user = ? AND notes.id = ?
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
