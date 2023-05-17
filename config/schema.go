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
	Status            int      `json:"status,omitempty"`
	Header            string   `json:"header,omitempty"`
	SkipHeaderValue   string   `json:"skip_header_value,omitempty"`
	ShouldHeaderValue string   `json:"should_header_value,omitempty"`
	Duration          Duration `json:"duration,omitempty"`
}

type Rule struct {
	Kind string `json:"kind"`
	Key  string `json:"key,omitempty"`

	Methods []string `json:"methods,omitempty"`
	Storage string   `json:"storage,omitempty"`

	TTL       string    `json:"ttl,omitempty"`
	EarlyTTL  string    `json:"early_ttl,omitempty"`
	Condition Condition `json:"condition,omitempty"`

	Count   int    `json:"count,omitempty"`
	Backoff string `json:"backoff,omitempty"`

	Hits            int `json:"hits,omitempty"`
	UpdateAfterHits int `json:"update_after,omitempty"`

	MinCalls  int    `json:"min_calls,omitempty"`
	ErrorRate int    `json:"error_rate,omitempty"`
	OpenTTL   string `json:"open_ttl,omitempty"`

	HeaderName string `json:"header,omitempty"`

	CheckReturnRequestId bool `json:"check_return,omitempty"`

	Response string `json:"response,omitempty"`

	Target string `json:"target,omitempty"`
	Path   string `json:"path,omitempty"`
	Route  string
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
