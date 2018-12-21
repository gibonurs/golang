package backup_test

import (
	"testing"

	"../backup"
	"github.com/stretchr/testify/require"
)

func TestDirHash(t *testing.T) {

	hash1a, err := backup.DirHash("test/hash1")
	require.NoError(t, err)
	hash1b, err := backup.DirHash("test/hash1")
	require.NoError(t, err)

	require.Equal(t, hash1a, hash1b, "hash1 i hash1b powinny być identyczne")

	hash2, err := backup.DirHash("test/hash2")
	require.NoError(t, err)

	require.NotEqual(t, hash1a, hash2, "hash1 i hash2 nie powinny być takie same")

}
