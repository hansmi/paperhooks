package client

import (
	"encoding/json"
	"fmt"
)

type Color struct {
	R, G, B uint8
}

var _ json.Marshaler = (*Color)(nil)
var _ json.Unmarshaler = (*Color)(nil)

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
