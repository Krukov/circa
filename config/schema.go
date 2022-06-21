package config

import (
	"encoding/json"
)

type GlobalOptions struct {
	Target         string   `json:"host"`
	Timeout        Duration `json:"timeout"`
	DefaultStorage string   `json:"default_storage"`
}

type Rule struct {
	Key  string `json:"key"`
	Kind string `json:"kind"`

	Timeout string   `json:"timeout"`
	Methods []string `json:"methods"`
	Storage string   `json:"storage"`

	TTL          string `json:"ttl"`
	EarlyTTL     string `json:"early_ttl"`
	CacheControl string `json:"cache_control"`

	Count   int    `json:"count"`
	Backoff string `json:"backoff"`

	Hits            int `json:"hits"`
	UpdateAfterHits int `json:"update_after"`

	RequestIDHeaderName string `json:"id_header"`

	SkipReturnRequestId bool `json:"skip_return"`

	Target string `json:"target"`
	Path   string `json:"path"`
	Method string `json:"method"`
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
