package resolver

import (
	"circa/rules"
	"errors"
	"math/rand"
	"sync"
)

var ErrNotFoundRule = errors.New("noRule")

type Resolver struct {
	router *node
	rules  map[string]*rules.Rule
	lock   *sync.RWMutex
}

func NewResolver() *Resolver {
	return &Resolver{
		router: newTrie(""),
		rules:  map[string]*rules.Rule{},
		lock:   &sync.RWMutex{},
	}
}

func (r *Resolver) Add(path string, rule *rules.Rule) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	ruleId := randStringBytes(36)
	r.router.addRule(path, ruleId)
	r.rules[ruleId] = rule
	return nil
}

func (r *Resolver) Resolve(path string) (rules []*rules.Rule, params map[string]string, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	var rules_ []string
	rules_, params, err = r.router.resolve(path)
	for _, _rule := range rules_ {
		rule, ok := r.rules[_rule]
		if ok {
			rules = append(rules, rule)
		} else {
			err = ErrNotFoundRule
		}
	}
	return rules, params, err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
