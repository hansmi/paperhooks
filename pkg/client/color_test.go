package client

import (
	"encoding/json"
	"image/color"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewColor(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input color.Color
		want  Color
	}{
		{
			name:  "black",
			input: color.Black,
		},
		{
			name:  "white",
			input: color.White,
			want:  Color{0xff, 0xff, 0xff},
		},
		{
			name:  "transparent",
			input: color.Transparent,
		},
		{
			name:  "opaque",
			input: color.Opaque,
			want:  Color{0xff, 0xff, 0xff},
		},
		{
			name:  "red from rgba",
			input: color.RGBA{R: 0xff, A: 0xff},
			want:  Color{0xff, 0, 0},
		},
		{
			name:  "green from nrgba",
			input: color.RGBA{G: 0xff, A: 0xff},
			want:  Color{0, 0xff, 0},
		},
		{
			name:  "semi-opaque blue from rgba",
			input: color.RGBA{B: 0x22, A: 0x7f},
			want:  Color{0, 0, 0x44},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := NewColor(tc.input)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Value diff (-want +got):\n%s", diff)
			}

			roundtrip := NewColor(got)

			if diff := cmp.Diff(tc.want, roundtrip); diff != "" {
				t.Errorf("Roundtrip diff (-want +got):\n%s", diff)
			}
		})
	}
}

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

func TestColorRGBA(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input Color
		want  color.Color
	}{
		{name: "black", want: color.Black},
		{name: "white", input: Color{0xff, 0xff, 0xff}, want: color.White},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := color.NRGBAModel

			if diff := cmp.Diff(model.Convert(tc.want), model.Convert(tc.input)); diff != "" {
				t.Errorf("Value diff (-want +got):\n%s", diff)
			}
		})
	}
}
