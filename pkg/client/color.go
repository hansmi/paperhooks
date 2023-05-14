package client

import (
	"encoding/json"
	"fmt"
	"image/color"
)

type Color struct {
	R, G, B uint8
}

var _ json.Marshaler = (*Color)(nil)
var _ json.Unmarshaler = (*Color)(nil)
var _ color.Color = (*Color)(nil)

// NewColor converts any color implementing the [color.Color] interface.
func NewColor(src color.Color) Color {
	nrgb := color.NRGBAModel.Convert(src).(color.NRGBA)

	return Color{
		R: nrgb.R,
		G: nrgb.G,
		B: nrgb.B,
	}
}

func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
}

func (c *Color) UnmarshalJSON(data []byte) error {
	var str *string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == nil {
		*c = Color{}
		return nil
	}

	if len(*str) == 7 && (*str)[0] == '#' {
		var r, g, b uint8

		if n, err := fmt.Sscanf(*str, "#%02x%02x%02x", &r, &g, &b); err == nil && n == 3 {
			c.R = r
			c.G = g
			c.B = b
			return nil
		}
	}

	return fmt.Errorf("unrecognized color format: %s", *str)
}

// RGBA implements [color.Color.RGBA].
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = 0xFFFF
	return
}
