package common

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

// RecordType is the type of record conveyed
type RecordType string

const (
	// AccountRecord is the Account record type
	AccountRecord RecordType = "account"
)

// Accounts holds list of accounts indexed by id
type Accounts map[int64]Account

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
