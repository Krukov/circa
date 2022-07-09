package resolver

import (
	"circa/rules"
	"errors"
	"sync"
)

var ErrNotFoundRule = errors.New("noRule")

type Resolver struct {
	router *node
	rules  map[string][]*rules.Rule
	lock   *sync.RWMutex
}

func NewResolver() *Resolver {
	return &Resolver{
		router: newTrie(""),
		rules:  map[string][]*rules.Rule{},
		lock:   &sync.RWMutex{},
	}
}

func (r *Resolver) Add(rule *rules.Rule) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.router.addRule(rule.Route)
	ruls, ok := r.rules[rule.Route]
	if !ok {
		ruls = []*rules.Rule{}
	}
	r.rules[rule.Route] = append(ruls, rule)
	return nil
}

func (r *Resolver) Resolve(path string) (rs []*rules.Rule, params map[string]string, err error) {
	var routes []string
	r.lock.RLock()
	defer r.lock.RUnlock()
	routes, params, err = r.router.resolve(path)
	for _, route := range routes {
		rs = append(rs, r.rules[route]...)
	}
	return rs, params, err
}
