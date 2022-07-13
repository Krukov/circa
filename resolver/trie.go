package resolver

import (
	"errors"
	"strings"
)

type node struct {
	name         string           // url prefix
	children     map[string]*node // {"users": n1, "products": n2}
	placeholders map[string]*node // ["{id}": n1, "{name}": {}]
	params       []string
	match        string // match rule for *
	rule         string
}

const STAR = "*"

var ErrNotFound = errors.New("noRoute")

func newTrie(name string) *node {
	return &node{name: name, children: map[string]*node{}, params: []string{}, placeholders: map[string]*node{}}
}

func (t *node) addRule(rulePath string) {
	sPath := strings.Split(strings.Trim(rulePath, "/"), "/")
	t.addNodes(sPath, rulePath)
}

func (t *node) addNodes(prefixs []string, rule string) {
	if len(prefixs) == 0 {
		return
	}
	var n *node
	if prefixs[0] == STAR {
		t.match = rule
		n = t
	} else {
		n = t.addNode(prefixs[0])
	}

	if len(prefixs) == 1 {
		n.rule = rule
		return
	}
	n.addNodes(prefixs[1:], rule)
}

// Add children to node with prefix
func (t *node) addNode(prefix string) *node {
	param := getParam(prefix) // like {id} -> id ; or test -> ""
	if param == "" {
		children, ok := t.children[prefix]
		if !ok {
			children = newTrie(prefix)
			t.children[prefix] = children
		}
		return children
	} else {
		placeholder, ok := t.placeholders[param]
		if !ok {
			placeholder = newTrie(prefix)
			placeholder.params = append(placeholder.params, param)
			t.placeholders[param] = placeholder
		}
		return placeholder
	}
}

func getParam(prefix string) string {
	if len(prefix) == 0 {
		return ""
	}
	if prefix[0] == '{' && prefix[len(prefix)-1] == '}' {
		return prefix[1 : len(prefix)-1]
	}
	return ""
}

func (t *node) resolve(path string) (names []string, params map[string]string, err error) {
	path = strings.Trim(strings.Trim(path, " "), "/")
	names, params, _ = t.walking(strings.Split(path, "/"), names, map[string]string{})
	if len(names) == 0 {
		err = ErrNotFound
		return
	}
	return
}

func (t *node) walking(path []string, rules []string, params map[string]string) ([]string, map[string]string, error) {
	var err error
	nameToFind := path[0]
	if nameToFind == "" {
		return rules, params, err
	}
	if t.match != "" {
		rules = append(rules, t.match)
	}

	if children, ok := t.children[nameToFind]; ok {
		if len(path) == 1 {
			if children.rule != "" && children.match == "" {
				rules = append(rules, children.rule)
			}
			return rules, params, nil
		}
		rules, params, err = children.walking(path[1:], rules, params)
		if err == nil {
			return rules, params, err
		}
	}
	for param, placeholder := range t.placeholders {
		params[param] = nameToFind
		if len(path) == 1 && placeholder.rule != "" {
			rules = append(rules, placeholder.rule)
			return rules, params, nil
		}
		rules, params, err = placeholder.walking(path[1:], rules, params)
		if err != nil {
			continue
		}
		return rules, params, err
	}
	return rules, params, ErrNotFound
}
