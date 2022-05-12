package common

import "encoding/json"

// Opaque is the type fit for "sub-record" of the record
type Opaque interface {
	Pack() (string, error)
}

// Pack converts Account to string
func (a Account) Pack() (string, error) {
	opaque, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(opaque), nil
}

// Pack converts Note to string
func (n Note) Pack() (string, error) {
	opaque, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(opaque), nil
}

// Pack converts Account to string
func (c Card) Pack() (string, error) {
	opaque, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(opaque), nil
}

// Pack returns Data field of Binary
func (b Binary) Pack() (string, error) {
	return b.Data, nil
}
