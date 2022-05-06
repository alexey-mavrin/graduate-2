package config

import (
	"errors"
	"os"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/alexey-mavrin/graduate-2/internal/crypt"
)

func checkFileMode(file string) (bool, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	perms := fileInfo.Mode().Perm()
	if (perms & 0077) != 0 {
		return false, nil
	}
	return true, nil
}

const minPhraseLen = 10

// GetKey returns a key computed from the given file content.
// It checks proper file permission and can check secret phrase strength.
func GetKey(file string) (*common.Key, error) {
	modeOK, err := checkFileMode(file)
	if err != nil {
		return nil, err
	}
	if !modeOK {
		return nil, errors.New("key phrase file mode incorrect")
	}
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(buf) < minPhraseLen {
		return nil, errors.New("key phrase is too short")
	}
	key := crypt.MakeKey(string(buf))
	return &key, nil
}
