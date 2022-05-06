package crypt

import (
	"github.com/alexey-mavrin/graduate-2/internal/common"
)

// EncryptAccount encrypts sensitive fields in Account record
func EncryptAccount(key common.Key, a common.Account) (common.Account, error) {
	e := common.Account{
		Name:     a.Name,
		UserName: a.UserName,
		URL:      a.URL,
	}
	ePassword, err := EncryptString(key, a.Password)
	if err != nil {
		return e, err
	}
	eMeta, err := EncryptString(key, a.Meta)
	if err != nil {
		return e, err
	}
	e.Password = ePassword
	e.Meta = eMeta
	return e, nil
}

// DecryptAccount decrypts sensitive fields in Account record
func DecryptAccount(key common.Key, e common.Account) (common.Account, error) {
	a := common.Account{
		Name:     e.Name,
		UserName: e.UserName,
		URL:      e.URL,
	}
	Password, err := DecryptString(key, e.Password)
	if err != nil {
		return e, err
	}
	Meta, err := DecryptString(key, e.Meta)
	if err != nil {
		return e, err
	}
	a.Password = Password
	a.Meta = Meta
	return a, nil
}

// EncryptNote encrypts sensitive fields in Note record
func EncryptNote(key common.Key, a common.Note) (common.Note, error) {
	e := common.Note{
		Name: a.Name,
	}
	eText, err := EncryptString(key, a.Text)
	if err != nil {
		return e, err
	}
	eMeta, err := EncryptString(key, a.Meta)
	if err != nil {
		return e, err
	}
	e.Text = eText
	e.Meta = eMeta
	return e, nil
}

// DecryptNote decrypts sensitive fields in Note record
func DecryptNote(key common.Key, e common.Note) (common.Note, error) {
	a := common.Note{
		Name: e.Name,
	}
	Text, err := DecryptString(key, e.Text)
	if err != nil {
		return e, err
	}
	Meta, err := DecryptString(key, e.Meta)
	if err != nil {
		return e, err
	}
	a.Text = Text
	a.Meta = Meta
	return a, nil
}
