package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"os"
	"sync"

	// sqlite sql package
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "secret_storage.db"
)

var secretStore *Store
var storeMutex sync.Mutex

// ErrNotFound is to indicate the absence of the record
var ErrNotFound = errors.New("Record not found")

// ErrAlreadyExists is to indicate the record already exist
var ErrAlreadyExists = errors.New("Entity already exists")

// Store is the secret storage
type Store struct {
	db *sql.DB
}

// User is the client of the secret store service
type User struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

// Account holds account data for some resource
type Account struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// Accounts holds list of accounts indexed by id
type Accounts map[int64]Account

// CloseDB closes database
func CloseDB() error {
	if secretStore != nil {
		err := secretStore.db.Close()
		if err != nil {
			return err
		}
		secretStore = nil
	}
	return nil
}

// DropStore removes the storage completely. Use for tests
func DropStore() error {
	err := CloseDB()
	if err != nil {
		return err
	}

	err = os.Remove(dbFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// NewStore initializes new storage or opens existing one
func NewStore() (*Store, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	if secretStore != nil {
		return secretStore, nil
	}

	secretStore = &Store{}
	var err error

	secretStore.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return secretStore, err
	}
	err = secretStore.db.Ping()
	if err != nil {
		return secretStore, err
	}

	_, err = secretStore.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		user TEXT NOT NULL UNIQUE CHECK (length(user) >= 3),
		full_name TEXT,
		password_hash TEXT
	)`)
	if err != nil {
		return secretStore, err
	}

	_, err = secretStore.db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL UNIQUE,
		url TEXT,
		user_name TEXT,
		password TEXT,
		FOREIGN KEY (user_id)
		  REFERENCES users (id)
		    ON DELETE CASCADE
		    ON UPDATE NO ACTION
	)`)
	if err != nil {
		return secretStore, err
	}

	return secretStore, nil
}

func (s *Store) isUserExists(userName string) (bool, error) {
	row := secretStore.db.QueryRow(
		`SELECT count(*) FROM users WHERE user = ?`,
		userName,
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

// CheckUserAuth checks the user password match
func (s *Store) CheckUserAuth(userName string, userPass string) (bool, error) {

	row := secretStore.db.QueryRow(
		`SELECT password_hash FROM users WHERE user = ?`,
		userName,
	)

	var dbPasswordHash string
	err := row.Scan(&dbPasswordHash)
	if err == sql.ErrNoRows {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256([]byte(userPass))
	passwordHash := hex.EncodeToString(hash[:])

	if dbPasswordHash == passwordHash {
		return true, nil
	}

	return false, nil
}

// AddUser creates user account
func (s *Store) AddUser(user User) (int64, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	hash := sha256.Sum256([]byte(user.Password))
	passwordHash := hex.EncodeToString(hash[:])

	userExists, err := s.isUserExists(user.Name)
	if err != nil {
		return 0, err
	}

	if userExists {
		return 0, ErrAlreadyExists
	}

	res, err := secretStore.db.Exec(`INSERT INTO users
		(user, full_name, password_hash)
		VALUES(?, ?, ?)`,
		user.Name,
		user.FullName,
		passwordHash,
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

// StoreAccount stores Account data for given user
func (s *Store) StoreAccount(user string, account Account) (int64, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`INSERT INTO accounts
		(user_id, name, url, user_name, password)
		VALUES((SELECT id from users where user=?), ?, ?, ?, ?)`,
		user,
		account.Name,
		account.URL,
		account.UserName,
		account.Password,
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
func (s *Store) UpdateAccount(user string, id int64, account Account) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	res, err := secretStore.db.Exec(`UPDATE accounts
		SET name = ?, url = ?, user_name = ?, password = ?
		WHERE id in
		( SELECT accounts.id FROM accounts
			JOIN users ON accounts.user_id = users.id
			WHERE users.user = ? AND accounts.id = ?)`,
		account.Name,
		account.URL,
		account.UserName,
		account.Password,
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

// GetAccounts returns list of stored accounts for the given user
// with all but password fileds filled
func (s *Store) GetAccounts(user string) (Accounts, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	accs := make(Accounts)
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
		var acc Account
		err = rows.Scan(&id, &acc.Name, &acc.URL, &acc.UserName)
		if err != nil {
			return accs, err
		}
		accs[id] = acc
	}
	return accs, nil
}

// GetAccount returns stored account record including the password
func (s *Store) GetAccount(user string, id int64) (Account, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	var acc Account

	row := secretStore.db.QueryRow(
		`SELECT accounts.name, accounts.url, accounts.user_name, accounts.password
			FROM accounts JOIN users ON accounts.user_id = users.id
			WHERE users.user = ? AND accounts.id = ?`,
		user, id,
	)

	err := row.Scan(&acc.Name, &acc.URL, &acc.UserName, &acc.Password)
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
