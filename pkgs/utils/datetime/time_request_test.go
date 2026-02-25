package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeRequest_ToTime(t *testing.T) {
	cases := []struct {
		name    string
		input   TimeRequest
		wantErr bool
	}{
		{name: "empty", input: "", wantErr: true},
		{name: "rfc3339", input: TimeRequest(time.Now().Format(time.RFC3339)), wantErr: false},
		{name: "date_time", input: "2025-01-02 15:04", wantErr: false},
		{name: "date_only", input: "2025-01-02", wantErr: false},
		{name: "invalid", input: "not-a-date", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.ToTime()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.False(t, got.IsZero())
		})
	}
}
