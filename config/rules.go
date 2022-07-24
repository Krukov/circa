package config

import (
	"circa/rules"
	"circa/storages"
	"fmt"
	"strings"
	"time"
)

func getRuleFromOptions(rule Rule, storage storages.Storage) (*rules.Rule, error) {
	processor, err := getRuleProcessorFromOptions(rule)
	if err != nil {
		return nil, err
	}
	methodsMap := map[string]bool{}
	if len(rule.Methods) > 0 {
		for _, method := range rule.Methods {
			methodsMap[strings.ToLower(method)] = true
		}
	} else {
		methodsMap["get"] = true
	}
	return &rules.Rule{
		Name:      rule.Kind,
		Key:       rule.Key,
		Route:     rule.Path,
		Methods:   methodsMap,
		Storage:   storage,
		Processor: processor,
		Condition: rules.Condition{
			Status:            rule.Condition.Status,
			SkipHeaderValue:   rule.Condition.SkipHeaderValue,
			ShouldHeaderValue: rule.Condition.ShouldHeaderValue,
			Header:            rule.Condition.Header,
			Duration:          rule.Condition.Duration.Duration,
		},
	}, nil
}

func getRuleProcessorFromOptions(rule Rule) (rules.RuleProcessor, error) {
	switch rule.Kind {
	case "proxy":
		return convertToProxyRule(rule)
	case "skip":
		return convertToSkipRule(rule)
	case "retry":
		return convertToRetryRule(rule)
	case "request-id":
		return convertToRequestIDRule(rule)
	case "proxy-header":
		return convertToProxyHeaderRule(rule)
	case "rate-limit":
		return convertToRateLimitRule(rule)
	case "idempotency":
		return convertToIdempotencyRule(rule)
	case "fail":
		return convertToFailRule(rule)
	case "hit":
		return convertToHitRule(rule)
	case "invalidate":
		return convertToInvalidateRule(rule)
	case "early":
		return convertToEarlyCacheRule(rule)
	case "cache":
		return convertToCacheRule(rule)
	case "glue":
		return convertToGlueRule(rule)
	case "":
		return convertToCacheRule(rule)
	}
	return nil, fmt.Errorf("unnown rule kind '%s'", rule.Kind)
}

func convertToProxyRule(rule Rule) (*rules.ProxyRule, error) {
	return &rules.ProxyRule{Target: rule.Target, Method: rule.Method, Path: rule.Path}, nil
}

func convertToSkipRule(rule Rule) (*rules.SkipRule, error) {
	return &rules.SkipRule{}, nil
}

func convertToRateLimitRule(rule Rule) (*rules.RateLimitRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.RateLimitRule{TTL: ttl, Limit: rule.Hits}, err
}

func convertToIdempotencyRule(rule Rule) (*rules.IdempotencyRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.IdempotencyRule{TTL: ttl}, err
}

func convertToRetryRule(rule Rule) (*rules.RetryRule, error) {
	backoff, err := timeFromString(rule.Backoff)
	return &rules.RetryRule{Count: rule.Count, Backoff: backoff}, err
}

func convertToInvalidateRule(rule Rule) (*rules.InvalidateRule, error) {
	methods := map[string]bool{}
	for _, method := range rule.Methods {
		methods[method] = true
	}
	return &rules.InvalidateRule{Methods: methods}, nil
}

func convertToCacheRule(rule Rule) (*rules.CacheRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.CacheRule{
		TTL:            ttl,
		Duration:       rule.Condition.Duration.Duration,
		ResponceStatus: rule.Condition.Status,
	}, err
}

func convertToEarlyCacheRule(rule Rule) (*rules.EarlyCacheRule, error) {
	ttl, err := timeFromString(rule.TTL)
	if err != nil {
		return nil, err
	}
	earlyTtl, errEarly := timeFromString(rule.EarlyTTL)
	return &rules.EarlyCacheRule{TTL: ttl, EarlyTTL: earlyTtl}, errEarly
}

func convertToFailRule(rule Rule) (*rules.FailRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.FailRule{TTL: ttl}, err
}

func convertToRequestIDRule(rule Rule) (*rules.RequestIDRule, error) {
	header := rule.HeaderName
	if header == "" {
		header = "X-Request-ID"
	}
	return &rules.RequestIDRule{HeaderName: header, SkipCheckReturn: rule.SkipReturnRequestId}, nil
}

func convertToProxyHeaderRule(rule Rule) (*rules.ProxyHeaderRule, error) {
	return &rules.ProxyHeaderRule{HeaderName: rule.HeaderName}, nil
}

func convertToHitRule(rule Rule) (*rules.HitRule, error) {
	ttl, err := timeFromString(rule.TTL)
	return &rules.HitRule{TTL: ttl, Hits: rule.Hits, UpdateAfterHits: rule.UpdateAfterHits}, err
}

func convertToGlueRule(rule Rule) (*rules.GlueRule, error) {
	return &rules.GlueRule{Calls: rule.Calls}, nil
}

func timeFromString(in string) (time.Duration, error) {
	if in == "" {
		return time.Second, nil
	}
	return time.ParseDuration(in)
}
