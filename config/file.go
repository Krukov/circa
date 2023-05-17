package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
	"time"
)

type fileConfigRepository struct {
	config *config
	path   string
}

func (repo *fileConfigRepository) GetStorages() (map[string]string, error) {
	return repo.config.Storages, nil
}

func (repo *fileConfigRepository) GetStorage(name string) (string, error) {
	for s, DSN := range repo.config.Storages {
		if s == name {
			return DSN, nil
		}
	}
	return "", errors.New("no storage")
}

func (repo *fileConfigRepository) AddStorage(name string, DSN string) error {
	repo.config.Storages[name] = DSN
	return nil
}

func (repo *fileConfigRepository) RemoveStorage(name string) error {
	if _, ok := repo.config.Storages[name]; !ok {
		return errors.New("no storage")
	}
	delete(repo.config.Storages, name)
	return nil
}

func (repo *fileConfigRepository) SetDefaultStorage(name string) error {
	repo.config.Options.DefaultStorage = name
	return nil
}

func (repo *fileConfigRepository) GetDefaultStorage() (string, error) {
	return repo.config.Options.DefaultStorage, nil
}

func (repo *fileConfigRepository) GetTarget() (string, error) {
	return repo.config.Options.Target, nil
}

func (repo *fileConfigRepository) SetTarget(host string) error {
	repo.config.Options.Target = host
	return nil
}

func (repo *fileConfigRepository) GetTimeout() (time.Duration, error) {
	return repo.config.Options.Timeout.Duration, nil
}

func (repo *fileConfigRepository) SetTimeout(timeout time.Duration) error {
	repo.config.Options.Timeout.Duration = timeout
	return nil
}

func (repo *fileConfigRepository) GetRoutes() ([]string, error) {
	routes := make([]string, len(repo.config.Rules))
	i := 0
	for route := range repo.config.Rules {
		routes[i] = route
		i += 1
	}
	sort.Strings(routes)
	return routes, nil
}

func (repo *fileConfigRepository) GetRules(route string) ([]Rule, error) {
	return repo.config.Rules[route], nil
}

func (repo *fileConfigRepository) AddRule(rule Rule) error {
	rules, ok := repo.config.Rules[rule.Route]
	if !ok {
		rules = []Rule{}
	}
	rules = append(rules, rule)
	repo.config.Rules[rule.Path] = rules
	return nil
}

func (repo *fileConfigRepository) RemoveRule(route, kind, key string) error {
	rules := []Rule{}
	for _, rule := range repo.config.Rules[route] {
		if rule.Kind == kind && rule.Key == key {
			continue
		}
		rules = append(rules, rule)
	}
	repo.config.Rules[route] = rules
	return nil
}

func (repo *fileConfigRepository) Sync() error {
	jsonRaw, err := json.MarshalIndent(repo.config, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(repo.path, jsonRaw, 0644)
}

func newFileConfig(path string) (configRepository, error) {
	jsonRaw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c, err := newFromJson(jsonRaw)
	if err != nil {
		return nil, err
	}
	conf := &fileConfigRepository{
		config: c,
		path:   path,
	}
	if c.Options.DefaultStorage == "" {
		for name := range c.Storages {
			conf.SetDefaultStorage(name)
			break
		}
	}
	return conf, nil
}
