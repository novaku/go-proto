package cache

import "encoding/json"

// ValueCodec serializes values for cache storage (JSON by default; swap for protobuf, etc.).
type ValueCodec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, dest any) error
}

// JSONCodec implements ValueCodec using encoding/json.
type JSONCodec struct{}

// Marshal encodes v as JSON.
func (JSONCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes JSON into dest.
func (JSONCodec) Unmarshal(data []byte, dest any) error {
	return json.Unmarshal(data, dest)
}
