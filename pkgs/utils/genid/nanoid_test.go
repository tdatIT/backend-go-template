package genid

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateNanoID_LengthAndAlphabet(t *testing.T) {
	id := GenerateNanoID()
	require.Len(t, id, nanoIDLength)

	for _, ch := range id {
		require.True(t, strings.ContainsRune(alphabet, ch))
	}
}
