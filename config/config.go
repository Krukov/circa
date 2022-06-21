package config

import (
	"errors"
	"sync"
	"time"

	"circa/resolver"
	"circa/rules"
	"circa/storages"

	"github.com/rs/zerolog/log"
)

type configRepository interface {
	GetStorages() (map[string]string, error)
	GetStorage(name string) (string, error)
	AddStorage(name string, DSN string) error
	RemoveStorage(name string) error

	SetDefaultStorage(name string) error
	GetDefaultStorage() (string, error)

	GetTarget() (string, error)

	GetTimeout() (time.Duration, error)

	GetRoutes() ([]string, error)
	GetRules(route string) ([]Rule, error)
	// AddRule(route string, rule Rule) error
	// RemoveRule(route, kind, key string) error

	// Sync()
}

type Config struct {
	repository configRepository
	storages   map[string]storages.Storage
	resolver   *resolver.Resolver

	target  string
	timeout time.Duration

	lock *sync.RWMutex
}

func (c *Config) Resolve(path string) (rules []*rules.Rule, params map[string]string, err error) {
	return c.resolver.Resolve(path)
}

func (c *Config) GetTarget() (string, error) {
	c.lock.RLock()
	if c.target != "" {
		defer c.lock.RUnlock()
		return c.target, nil
	}
	c.lock.RUnlock()
	c.lock.Lock()
	defer c.lock.Unlock()
	target, err := c.repository.GetTarget()
	if err != nil {
		return "", err
	}
	c.target = target
	return target, nil
}

func (c *Config) GetTimeout() (time.Duration, error) {
	c.lock.RLock()
	if c.timeout != 0 {
		defer c.lock.RUnlock()
		return c.timeout, nil
	}
	c.lock.RUnlock()
	// Take write lock
	c.lock.Lock()
	defer c.lock.Unlock()
	timeout, err := c.repository.GetTimeout()
	if err != nil {
		return 0, err
	}
	c.timeout = timeout
	return timeout, nil
}

func NewConfigFromDSN(dsn string, resolver *resolver.Resolver) (*Config, error) {
	repo, err := newFileConfig(dsn)
	if err != nil {
		return nil, err
	}
	storagesMap := map[string]storages.Storage{}
	storagesFromConfig, err := repo.GetStorages()
	if err != nil {
		return nil, err
	}

	for name, DSN := range storagesFromConfig {
		storagesMap[name], err = storages.StorageFormDSN(DSN)
		if err != nil {
			return nil, err
		}
		log.Info().Msgf("Configured storage '%v' with dns '%v'", name, DSN)
	}

	// defRequest := &message.Request{Host: c.Options.Target, Timeout: timeout} // Move to runner
	routes, err := repo.GetRoutes()
	if err != nil {
		return nil, err
	}
	var storage storages.Storage
	var ok bool
	defaultStorageName, err := repo.GetDefaultStorage()
	if err != nil {
		return nil, err
	}
	defaultStorage, ok := storagesMap[defaultStorageName]
	if !ok {
		return nil, errors.New("wrong default Storage setup")
	}

	for _, route := range routes {
		rules, err := repo.GetRules(route)
		if err != nil {
			return nil, err
		}
		for _, ruleOptions := range rules {
			storage, ok = storagesMap[ruleOptions.Storage]
			if !ok {
				storage = defaultStorage
			}
			rule, err := getRuleFromOptions(ruleOptions, storage, route)
			if err != nil {
				return nil, err
			}

			resolver.Add(route, rule)
		}
	}
	return &Config{
		repository: repo,
		resolver:   resolver,
		storages:   storagesMap,
		lock:       &sync.RWMutex{},
	}, nil
}
