package datetime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTimeNowWithUTC7_Location(t *testing.T) {
	localTime := GetTimeNowWithUTC7()
	require.Equal(t, "Asia/Bangkok", localTime.Location().String())
}
