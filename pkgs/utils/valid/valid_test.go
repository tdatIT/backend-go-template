package valid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type sampleRequest struct {
	Name string `validate:"required"`
}

func TestGetValidator_Singleton(t *testing.T) {
	v1 := GetValidator()
	v2 := GetValidator()
	require.Same(t, v1, v2)
}

func TestValidator_Validate(t *testing.T) {
	v := GetValidator()

	err := v.Validate(sampleRequest{})
	require.Error(t, err)

	err = v.Validate(sampleRequest{Name: "ok"})
	require.NoError(t, err)
}
