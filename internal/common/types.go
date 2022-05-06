package common

// Key is the AES key type used
type Key [32]byte

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
	Meta     string `json:"meta"`
}

// Note holds text data
type Note struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Meta string `json:"meta"`
}

// Record can hold any record that could be stored
type Record struct {
	Type    RecordType
	Account *Account
	Note    *Note
}

// Records can hold the map of any record that could be stored
type Records struct {
	Type     RecordType
	Accounts *Accounts
	Notes    *Notes
}

// RecordType is the type of record conveyed
type RecordType string

const (
	// AccountRecord is the Account record type
	AccountRecord RecordType = "account"
	// NoteRecord is the Note record type
	NoteRecord RecordType = "note"
)

// Accounts holds list of accounts indexed by id
type Accounts map[int64]Account

// Notes holds list of notes indexed by id
type Notes map[int64]Note

// AddUserResponse is the response for AddUser request
type AddUserResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

// StoreRecordResponse is the responce for store account
type StoreRecordResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	ID     int64  `json:"id"`
}
