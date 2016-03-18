//go:generate ffjson $GOFILE
package eventsourcedb

import "encoding/json"

type Event struct {
	ID     uint64          `json:"id"`
	Stream string          `json:"stream"`
	Type   string          `json:"type"`
	Body   json.RawMessage `json:"body"`
}
