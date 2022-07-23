package config

import (
	"encoding/json"
)

type GlobalOptions struct {
	Target         string   `json:"target"`
	Timeout        Duration `json:"timeout,omitempty"`
	DefaultStorage string   `json:"default_storage"`
}

type Condition struct {
	Status      int      `json:"status,omitempty"`
	Header      string   `json:"header,omitempty"`
	HeaderValue string   `json:"header_value,omitempty"`
	Duration    Duration `json:"duration,omitempty"`
}

type Rule struct {
	Kind string `json:"kind"`
	Key  string `json:"key,omitempty"`

	Methods []string `json:"methods,omitempty"`
	Storage string   `json:"storage,omitempty"`

	TTL          string    `json:"ttl,omitempty"`
	EarlyTTL     string    `json:"early_ttl,omitempty"`
	Condition    Condition `json:"condition,omitempty"`

	Count   int    `json:"count,omitempty"`
	Backoff string `json:"backoff,omitempty"`

	Hits            int `json:"hits,omitempty"`
	UpdateAfterHits int `json:"update_after,omitempty"`

	RequestIDHeaderName string `json:"id_header,omitempty"`

	SkipReturnRequestId bool `json:"skip_return,omitempty"`

	Target string `json:"target,omitempty"`
	Path   string `json:"path,omitempty"`
	Method string `json:"method,omitempty"`

	Calls map[string]string `json:"calls,omitempty"`
}

type config struct {
	Version  string            `json:"version"`
	Storages map[string]string `json:"storages"`
	Rules    map[string][]Rule `json:"rules"`
	Options  GlobalOptions     `json:"options"`
}

func newFromJson(in []byte) (*config, error) {
	c := &config{}
	if err := json.Unmarshal(in, &c); err != nil {
		return nil, err
	}
	return c, nil
}
