package config

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/alexey-mavrin/graduate-2/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTmpFile(t *testing.T, str string, mode fs.FileMode) string {
	tmpDir := os.TempDir()
	tmpFile := path.Join(tmpDir, "key.txt")
	err := os.WriteFile(tmpFile, []byte(str), mode)
	require.NoError(t, err)
	return tmpFile
}

func Test_GetKey(t *testing.T) {
	tests := []struct {
		name         string
		fileMode     fs.FileMode
		secretPhrase string
		wantErr      assert.ErrorAssertionFunc
		wantNilKey   bool
	}{
		{
			name:         "normal run",
			fileMode:     0600,
			secretPhrase: "1234567890",
			wantErr:      assert.NoError,
		},
		{
			name:         "wrong mode",
			fileMode:     0640,
			secretPhrase: "1234567890",
			wantErr:      assert.Error,
			wantNilKey:   true,
		},
		{
			name:         "phrase is too short",
			fileMode:     0600,
			secretPhrase: "123456789",
			wantErr:      assert.Error,
			wantNilKey:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := writeTmpFile(t, tt.secretPhrase, tt.fileMode)
			key, err := GetKey(file)
			os.Remove(file)
			tt.wantErr(t, err)
			if tt.wantNilKey {
				if key != nil {
					t.Error("key is not nil as required")
				}
			} else {
				var k common.Key
				if len(key) != len(k) {
					t.Error("wrong key size")
				}
			}
		})
	}
}
