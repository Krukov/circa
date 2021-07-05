package handler

import (
	"errors"
	"strings"
)

type ruleName string

type node struct {
	children map[string]*node // {"users": n, "products": n2}
	params   []string
	ruleName ruleName
}

var NotFound = errors.New("no route")
const STAR = "*"

func newTrie() *node {
	return &node{children: map[string]*node{}, params: []string{}}
}

func (t *node) addRule (rulePath string, ruleName ruleName) {
	rulePath = strings.Trim(rulePath, "/")
	sPath := strings.Split(rulePath, "/")
	t.addPrefixs(sPath, ruleName)
}


func (t *node) removeRule (ruleName ruleName) {
	for name, children := range t.children {
		if children.ruleName == ruleName {
			delete(t.children, name)
		}
		if len(children.children) > 0 {
			children.removeRule(ruleName)
		}
	}
}

func (t *node) addPrefixs (prefixs []string, ruleName ruleName) {
	if len(prefixs) == 0 {
		return
	}
	n := t.addPrefix(prefixs[0])
	if len(prefixs) == 1 {
		n.ruleName = ruleName
		return
	}
	n.addPrefixs(prefixs[1:], ruleName)
}

func (t *node) addPrefix (prefix string) *node {
	name, params := getNameStar(prefix)
	children, ok := t.children[name]
	if !ok {
		children = &node{children: map[string]*node{}, params: []string{}}
		t.children[name] = children
	}
	for _, param := range params {
		t.params = append(t.children[name].params, param)
		t.children[name].params = append(t.children[name].params, param)
	}
	return children
}

func getNameStar(prefix string) (string, []string) {
	if len(prefix) == 0 {
		return prefix, []string{}
	}
	if prefix[0] == '{' && prefix[len(prefix) - 1] == '}' {
		return STAR, []string{prefix[1:len(prefix) - 1]}
	}
	return prefix, []string{}
}

func (t *node) getRoute (path string) (name ruleName, params map[string]string, err error) {
	var n *node
	path = strings.Trim(path, "/")
	sPath := strings.Split(path, "/")
	n, params, err = t.getNode(sPath, map[string]string{})
	if n != nil {
		name = n.ruleName
	}
	if name == "" {
		err = NotFound
	}
	return
}

func (t *node) getNode (path []string, params map[string]string) (*node, map[string]string, error) {
	nameToFind := path[0]

	for name, children := range t.children {
		if name == nameToFind {
			if len(path) == 1 {
				return children, params, nil
			}
			resNode, resParams, err := children.getNode(path[1:], params)
			if err == nil {
				return resNode, resParams, err
			}
		}
	}
	if n, ok := t.children[STAR]; ok {
		for _, param := range t.params {
			params[param] = nameToFind
		}
		if len(path) == 1 {
			return n, params, nil
		}
		return n.getNode(path[1:], params)
	}

	return nil, params, NotFound
}