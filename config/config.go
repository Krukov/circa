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
	AddStorage(name, DSN string) error
	RemoveStorage(name string) error

	SetDefaultStorage(name string) error
	GetDefaultStorage() (string, error)

	GetTarget() (string, error)

	GetTimeout() (time.Duration, error)

	GetRoutes() ([]string, error)
	GetRules(route string) ([]Rule, error)
	AddRule(rule Rule) error
	// RemoveRule(route, kind, key string) error

	Sync() error
}

type Config struct {
	repository  configRepository
	storages    map[string]storages.Storage
	storageLock *sync.RWMutex

	resolver *resolver.Resolver

	target  string
	timeout time.Duration

	lock *sync.RWMutex
}

func (c *Config) Init() error {
	rules, err := c.getRules()
	if err != nil {
		return err
	}
	for _, rule := range rules {
		c.resolver.Add(rule)
	}
	return nil
}

func (c *Config) AddRule(rule Rule) error {
	var err error
	c.storageLock.RLock()
	storageClass, ok := c.storages[rule.Storage]
	if !ok {
		storageClass, _, err = c.getDefaultStorage()
		if err != nil {
			c.storageLock.RUnlock()
			return err
		}
	}
	c.storageLock.RUnlock()
	ruleRule, err := getRuleFromOptions(rule, storageClass)
	if err != nil {
		return err
	}
	c.resolver.Add(ruleRule)
	err = c.repository.AddRule(rule)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) getRules() ([]*rules.Rule, error) {
	routes, err := c.repository.GetRoutes()
	if err != nil {
		return nil, err
	}
	var storage storages.Storage
	var ok bool
	c.storageLock.RLock()
	defaultStorage, defaultStorageName, err := c.getDefaultStorage()
	c.storageLock.RUnlock()
	if err != nil {
		return nil, err
	}
	returnRules := []*rules.Rule{}

	for _, route := range routes {
		rules, err := c.repository.GetRules(route)
		if err != nil {
			return nil, err
		}
		for _, ruleOptions := range rules {
			storageName := ruleOptions.Storage
			c.storageLock.RLock()
			storage, ok = c.storages[ruleOptions.Storage]
			c.storageLock.RUnlock()
			if !ok {
				storage = defaultStorage
				storageName = defaultStorageName
			}
			if ruleOptions.Path == "" {
				ruleOptions.Path = route
			}
			ruleOptions.Route = route
			rule, err := getRuleFromOptions(ruleOptions, storage)
			if err != nil {
				return nil, err
			}
			rule.StorageName = storageName
			returnRules = append(returnRules, rule)
		}
	}
	return returnRules, nil
}

func (c *Config) getDefaultStorage() (storages.Storage, string, error) {
	defaultStorageName, err := c.repository.GetDefaultStorage()
	if err != nil {
		return nil, "", err
	}
	defaultStorage, ok := c.storages[defaultStorageName]
	if !ok {
		return nil, "", errors.New("wrong default Storage setup")
	}
	return defaultStorage, defaultStorageName, nil
}

func (c *Config) Resolve(path string) (rules []*rules.Rule, params map[string]string, err error) {
	rules, params, err = c.resolver.Resolve(path)

	// TODO: redone without lock storages check (for many rules we don't need storage at all)
	c.storageLock.RLock()
	defaultStorage, defaultStorageName, _ := c.getDefaultStorage()

	// Rewrite storage if it was deleted
	for _, r := range rules {
		_, ok := c.storages[r.StorageName]
		if !ok {
			r.Storage = defaultStorage
			r.StorageName = defaultStorageName
		}
	}
	c.storageLock.RUnlock()
	return rules, params, err
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

func (c *Config) GetStorages() (map[string]string, error) {
	return c.repository.GetStorages()
}

func (c *Config) AddStorage(name, DSN string) error {
	var err error
	c.storageLock.Lock()
	c.storages[name], err = storages.StorageFormDSN(DSN)
	if err != nil {
		return err
	}
	c.storageLock.Unlock()
	return c.repository.AddStorage(name, DSN)
}

func (c *Config) DelStorage(name string) error {
	err := c.repository.RemoveStorage(name)
	if err != nil {
		return err
	}
	c.storageLock.Lock()
	delete(c.storages, name)
	c.storageLock.Unlock()
	return nil
}

func (c *Config) GetRoutes() ([]*rules.Rule, error) {
	return c.getRules()
}

func (c *Config) Sync() error {
	return c.repository.Sync()
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
	c := Config{
		repository:  repo,
		resolver:    resolver,
		storages:    storagesMap,
		storageLock: &sync.RWMutex{},
		lock:        &sync.RWMutex{},
	}
	return &c, c.Init()
}
