package client

import (
	"encoding/json"

	"golang.org/x/exp/maps"
)

type objectFields map[string]any

func (f objectFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(f))
}

// AsMap returns a map from object field name to value.
func (f objectFields) AsMap() map[string]any {
	return maps.Clone(f)
}

func (f objectFields) set(name string, value any) {
	f[name] = value
}

func (f objectFields) build() map[string]any {
	return f
}
