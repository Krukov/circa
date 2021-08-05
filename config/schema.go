package config

import "encoding/json"

type GlobalOptions struct {
	Target  string `json:"host"`
	Timeout string `json:"timeout"`
}

type Rule struct {
	Key  string `json:"key"`
	Type string `json:"type"`

	Timeout string   `json:"timeout"`
	Methods []string `json:"methods"`
	Storage string   `json:"storage"`

	TTL          string `json:"ttl"`
	EarlyTTL     string `json:"early_ttl"`
	CacheControl string `json:"cache_control"`

	Count            int `json:"count"`

	Hits int `json:"hits"`
	UpdateAfterHits int `json:"update_after"`

	RequestIDHeaderName string `json:"id_header"`
}

type config struct {
	Version  string                     `json:"version"`
	Storages map[string]string          `json:"storages"`
	Rules    map[string][]Rule          `json:"rules"`
	Options  GlobalOptions              `json:"options"`
}

func newFromJson(in []byte) (*config, error) {
	c := &config{}
	if err := json.Unmarshal(in, &c); err != nil {
		return nil, err
	}
	return c, nil
}
