package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFilter(t *testing.T) {
	tests := []struct {
		name          string
		args          string
		want          []Filter
		errorExpected bool
	}{
		{
			name:          "no filter",
			args:          "",
			want:          []Filter{},
			errorExpected: false,
		},
		{
			name: "lots of spaces",
			args: "userName              eq            \"test\"",
			want: []Filter{
				{
					FilterField:    Username,
					FilterOperator: Eq,
					FilterValue:    "test",
				},
			},
			errorExpected: false,
		},
		{
			name:          "invalid operator",
			args:          "userName foo \"test\"",
			want:          []Filter{},
			errorExpected: true,
		},
		{
			name:          "invalid field",
			args:          "user eq \"test\"",
			want:          []Filter{},
			errorExpected: true,
		},
		{
			name:          "malformed value",
			args:          "userName eq \"test",
			want:          []Filter{},
			errorExpected: true,
		},
		{
			name: "valid filter",
			args: "userName eq \"test\"",
			want: []Filter{
				{
					FilterField:    Username,
					FilterOperator: Eq,
					FilterValue:    "test",
				},
			},
			errorExpected: false,
		},
		{
			name: "valid filter with uppercase operator",
			args: "userName EQ \"test\"",
			want: []Filter{
				{
					FilterField:    Username,
					FilterOperator: Eq,
					FilterValue:    "test",
				},
			},
			errorExpected: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseFilter(tc.args)

			if tc.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
