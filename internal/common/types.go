package common

import "encoding/json"

const (
	EMPTY_STRING = ""
)

// RequestData struct
type RequestData struct {
	Service string          `json:"service"`
	Data    json.RawMessage `json:"data"`
}
