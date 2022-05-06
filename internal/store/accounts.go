package store

import (
	"database/sql"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	// sqlite sql package
	_ "github.com/mattn/go-sqlite3"
)

// StoreAccountID stores account with the ID specified
func (s *Store) StoreAccountID(id int64, user string, account common.Account) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	_, err := secretStore.db.Exec(`INSERT INTO accounts
		(id, user_id, name, url, user_name, password, meta)
		VALUES(?, (SELECT id from users where user=?), ?, ?, ?, ?, ?)`,
		id,
		user,
		account.Name,
		account.URL,
		account.UserName,
		account.Password,
		account.Meta,
	)
	if err != nil {
		return err
	}

	return nil
}

// StoreAccount stores Account data for given user
func (s *Store) StoreAccount(user string, account common.Account) (int64, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`INSERT INTO accounts
		(user_id, name, url, user_name, password, meta)
		VALUES((SELECT id from users where user=?), ?, ?, ?, ?, ?)`,
		user,
		account.Name,
		account.URL,
		account.UserName,
		account.Password,
		account.Meta,
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

// UpdateAccount updates an accout record
func (s *Store) UpdateAccount(user string, id int64, account common.Account) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`UPDATE accounts
		SET name = ?, url = ?, user_name = ?, password = ?, meta = ?
		WHERE id in
		( SELECT accounts.id FROM accounts
			JOIN users ON accounts.user_id = users.id
			WHERE users.user = ? AND accounts.id = ?)`,
		account.Name,
		account.URL,
		account.UserName,
		account.Password,
		account.Meta,
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

// ListAccounts returns list of stored accounts for the given user
// with all but password fileds filled
func (s *Store) ListAccounts(user string) (common.Accounts, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	accs := make(common.Accounts)
	rows, err := secretStore.db.Query(
		`SELECT accounts.id, accounts.name, accounts.url, accounts.user_name
			FROM accounts JOIN users ON accounts.user_id = users.id
			WHERE users.user = ?`,
		user,
	)
	if err != nil {
		return accs, err
	}

	for rows.Next() {
		var id int64
		var acc common.Account
		err = rows.Scan(&id, &acc.Name, &acc.URL, &acc.UserName)
		if err != nil {
			return accs, err
		}
		accs[id] = acc
	}
	return accs, nil
}

// GetAccount returns stored account record including the password
func (s *Store) GetAccount(user string, id int64) (common.Account, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	var acc common.Account

	row := secretStore.db.QueryRow(
		`SELECT accounts.name, accounts.url, accounts.user_name,
			accounts.password, accounts.meta
			FROM accounts JOIN users ON accounts.user_id = users.id
			WHERE users.user = ? AND accounts.id = ?`,
		user, id,
	)

	err := row.Scan(&acc.Name,
		&acc.URL,
		&acc.UserName,
		&acc.Password,
		&acc.Meta,
	)
	if err == sql.ErrNoRows {
		return acc, ErrNotFound
	}
	if err != nil {
		return acc, err
	}
	return acc, nil
}

// DeleteAccount deletes the specified account record
func (s *Store) DeleteAccount(user string, id int64) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(
		`DELETE FROM accounts
		 	WHERE id in
			( SELECT accounts.id FROM accounts
				JOIN users ON accounts.user_id = users.id
				WHERE users.user = ? AND accounts.id = ?
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
