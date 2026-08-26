package grafana

import "encoding/json"

type Policy struct {
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Body      json.RawMessage `json:"body"`
	ExpiresAt string          `json:"expires_at,omitempty"`
}

type StringAttribute struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type Drop struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}
