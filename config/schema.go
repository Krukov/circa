package config

import "encoding/json"

type GlobalOptions struct {
	Target  string `json:"host"`
	Timeout string `json:"timeout"`
}

type Rule struct {
	Type string `json:"type"`

	Timeout string   `json:"timeout"`
	Methods []string `json:"methods"`
	Storage string   `json:"storage"`

	TTL          string `json:"ttl"`
	EarlyTTL     string `json:"early_ttl"`
	CacheControl string `json:"cache_control"`

	Count int `json:"count"`
}

type config struct {
	Version  string                     `json:"version"`
	Storages map[string]string          `json:"storages"`
	Rules    map[string]map[string]Rule `json:"rules"`
	Options  GlobalOptions              `json:"options"`
}

func newFromJson(in []byte) (*config, error) {
	c := &config{}
	if err := json.Unmarshal(in, &c); err != nil {
		return nil, err
	}
	return c, nil
}
