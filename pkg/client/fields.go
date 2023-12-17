package client

import "encoding/json"

type objectFields map[string]any

func (f objectFields) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(f))
}

func (f objectFields) set(name string, value any) {
	f[name] = value
}

func (f objectFields) build() map[string]any {
	return f
}
