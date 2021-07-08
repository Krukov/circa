package handler

import (
	"errors"
	"strings"
)

type ruleName string

type node struct {
	name      string
	children  map[string]*node // {"users": n1, "products": n2}
	params    []string
	rule      ruleName
	downRules []ruleName
}

var NotFound = errors.New("no route")
var DownRuleError = errors.New("down Rule")

const STAR = "*"
const PLACEHOLDER = "."

func newTrie() *node {
	return &node{name: "root", children: map[string]*node{}, params: []string{}, downRules: []ruleName{}}
}

func (t *node) addRule(rulePath string, rule ruleName) {
	rulePath = strings.Trim(rulePath, "/")
	sPath := strings.Split(rulePath, "/")
	t.addPrefixs(sPath, rule)
}

func (t *node) setDownRule(rule ruleName) {
	t.downRules = append(t.downRules, rule)
	for _, c := range t.children {
		c.setDownRule(rule)
	}
}

func (t *node) addPrefixs(prefixs []string, rule ruleName) {
	if len(prefixs) == 0 {
		return
	}
	n := t.addPrefix(prefixs[0])

	if len(prefixs) == 1 {
		if n.name == STAR {
			t.setDownRule(rule)
		} else {
			n.rule = rule
		}
		return
	}
	n.addPrefixs(prefixs[1:], rule)
}

// Add children to node with prefix
func (t *node) addPrefix(prefix string) *node {
	name, param := getNameStar(prefix) // like {id} -> PLACEHOLDER, id ; or test -> test, nil
	children, ok := t.children[name]
	if !ok {
		children = &node{name: prefix, children: map[string]*node{}, params: []string{}, downRules: []ruleName{}}
		children.downRules = t.downRules
		t.children[name] = children
	}

	if param != "" {
		t.params = append(t.children[name].params, param)
		t.children[name].params = append(t.children[name].params, param)
	}
	return children
}

func getNameStar(prefix string) (string, string) {
	if len(prefix) == 0 {
		return prefix, ""
	}
	if prefix[0] == '{' && prefix[len(prefix)-1] == '}' {
		return PLACEHOLDER, prefix[1 : len(prefix)-1]
	}
	return prefix, ""
}

func (t *node) resolve(path string) (names []ruleName, params map[string]string, err error) {
	var n *node
	path = strings.Trim(path, "/")
	sPath := strings.Split(path, "/")
	n, params, err = t.getNode(sPath, map[string]string{})
	if n != nil {
		if err != DownRuleError && n.rule != "" {
			names = append(names, n.rule)
		}
		names = append(names, n.downRules...)
		err = nil
	}
	if len(names) == 0 {
		err = NotFound
	}
	return
}

func (t *node) getNode(path []string, params map[string]string) (*node, map[string]string, error) {
	nameToFind := path[0]

	if children, ok := t.children[nameToFind]; ok {
		if len(path) == 1 {
			return children, params, nil
		}
		resNode, resParams, err := children.getNode(path[1:], params)
		if err == nil {
			return resNode, resParams, err
		}
		if err == DownRuleError {
			return resNode, resParams, err
		}
	}

	if n, ok := t.children[PLACEHOLDER]; ok {
		for _, param := range t.params {
			params[param] = nameToFind
		}
		if len(path) == 1 {
			return n, params, nil
		}
		resNode, resParams, err := n.getNode(path[1:], params)

		if err == nil {
			return resNode, resParams, err
		}
		if err == DownRuleError {
			return resNode, resParams, err
		}
	}

	if len(t.downRules) > 0 {
		return t, params, DownRuleError
	}
	return nil, params, NotFound
}
