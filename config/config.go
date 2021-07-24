package config

import (
	"io/ioutil"
	"time"

	"circa/handler"
	"circa/message"
	"circa/rules"
	"circa/storages"

	"github.com/rs/zerolog/log"
)

func AdjustJsonConfig(r *handler.Runner, path string) error {
	jsonRaw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	c, err := newFromJson(jsonRaw)
	if err != nil {
		return err
	}

	storagesMap := map[string]storages.Storage{}
	var defaultStorage storages.Storage
	for name, DSN := range c.Storages {
		storagesMap[name], err = storages.StorageFormDSN(DSN)
		if err != nil {
			return err
		}
		log.Info().Msgf("Configured storage '%v' with dns '%v'", name, DSN)
		defaultStorage = storagesMap[name]
	}

	timeout, err := timeFromString(c.Options.Timeout)
	if err != nil {
		return err
	}
	defRequest := &message.Request{Host: c.Options.Target, Timeout: timeout}
	var rule rules.Rule
	var ok bool
	var storage storages.Storage
	for rulePath, rule_ := range c.Rules {
		for temp, ruleOptions := range rule_ {
			rule, err = getRuleFromOptions(ruleOptions)
			if err != nil {
				return err
			}
			storage, ok = storagesMap[ruleOptions.Storage]
			if !ok {
				storage = defaultStorage
			}
			r.AddHandlers(rulePath, handler.NewHandler(rule, storage, temp, defRequest, ruleOptions.Methods))
		}
	}
	r.SetProxy(c.Options.Target, timeout)
	return nil

}

func getRuleFromOptions(rule Rule) (rules.Rule, error) {
	switch rule.Type {
	case "proxy":
		return convertToProxyRule(rule)
	case "retry":
		return convertToRetryRule(rule)
	case "request_id":
		return convertToRequestIDRule(rule)
	case "fail":
		return convertToFailRule(rule)
	case "hit":
		return convertToHitRule(rule)
	case "simple":
	}
	return convertToCacheRule(rule)
}

func convertToProxyRule(rule Rule) (*rules.ProxyRule, error) {
	return &rules.ProxyRule{}, nil
}

func convertToRetryRule(rule Rule) (*rules.RetryRule, error) {
	return &rules.RetryRule{Count: rule.Count}, nil
}

func convertToCacheRule(rule Rule) (*rules.CacheRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.CacheRule{TTL: ttl}, err
}

func convertToFailRule(rule Rule) (*rules.FailRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.FailRule{TTL: ttl}, err
}

func convertToRequestIDRule(rule Rule) (*rules.RequestIDRule, error) {
	return &rules.RequestIDRule{HeaderName: rule.RequestIDHeaderName}, nil
}

func convertToHitRule(rule Rule) (*rules.HitRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.HitRule{TTL: ttl, Hits: rule.Hits, UpdateAfterHits: rule.UpdateAfterHits}, err
}

func timeFromString(in string) (time.Duration, error) {
	if in == "" {
		return time.Second, nil
	}
	return time.ParseDuration(in)
}
