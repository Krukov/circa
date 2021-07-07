package config

import (
	"circa/handler"
	"circa/message"
	"circa/rules"
	"circa/storages"
	"errors"
	"io/ioutil"
	"time"
)

func AdjustJsonConfig (r *handler.Runner, path string) error {
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
				if len(storagesMap) == 1 {
					storage = defaultStorage
				} else {
					return errors.New("unknown storage")
				}
			}
			r.AddHandlers(rulePath, handler.NewHandler(rule, storage, temp, defRequest, ruleOptions.Methods))
		}
	}
	r.SetProxy(c.Options.Target, timeout)
	return nil

}


func getRuleFromOptions (rule Rule) (rules.Rule, error) {
	switch rule.Type {
	case "proxy":
		return convertToProxyRule(rule)
	case "fail":
		return convertToFailRule(rule)
	case "simple":
	}
	return convertToCacheRule(rule)
}

func convertToProxyRule(rule Rule) (*rules.ProxyRule, error) {
	return &rules.ProxyRule{}, nil
}

func convertToCacheRule(rule Rule) (*rules.CacheRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.CacheRule{TTL: ttl}, err
}

func convertToFailRule(rule Rule) (*rules.FailRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.FailRule{TTL: ttl}, err
}


func timeFromString(in string) (time.Duration, error) {
	if in == "" {
		return time.Second, nil
	}
	return time.ParseDuration(in)
}