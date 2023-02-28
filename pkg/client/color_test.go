package client

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestColorJSON(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value Color
		want  string
	}{
		{
			name: "zero",
			want: `"#000000"`,
		},
		{
			name:  "white",
			value: Color{0xFF, 0xFF, 0xFF},
			want:  `"#ffffff"`,
		},
		{
			name:  "red",
			value: Color{R: 0xFF},
			want:  `"#ff0000"`,
		},
		{
			name:  "blue",
			value: Color{B: 0xFF},
			want:  `"#0000ff"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.value)
			if err != nil {
				t.Errorf("Marshalling failed: %v", err)
			}

			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Errorf("Marshalled value diff (-want +got):\n%s", diff)
			}

			var restored Color

			if err := json.Unmarshal(got, &restored); err != nil {
				t.Errorf("Unmarshalling failed: %v", err)
			}

			if diff := cmp.Diff(tc.value, restored); diff != "" {
				t.Errorf("Value diff (-want +got):\n%s", diff)
			}
		})
	}
}
