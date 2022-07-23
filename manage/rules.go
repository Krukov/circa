package manage

import (
	"circa/config"
	"circa/rules"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type Rule struct {
	Path    string   `json:"path"`
	Key     string   `json:"key"`
	Kind    string   `json:"kind"`
	Storage string   `json:"storage"`
	Methods []string `json:"methods"`
}

var rulePath = regexp.MustCompile("^/api/rules/(.+)$")

func (cm *configManage) Rules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	m := rulePath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		if r.Method == http.MethodGet {
			cm.GetRules(w, r)
			return
		}
		if r.Method == http.MethodPost {
			cm.AddRoute(w, r)
			return
		}
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
}

func (cm *configManage) GetRules(w http.ResponseWriter, r *http.Request) {
	var rs []*rules.Rule
	var err error
	if r.URL.Query().Get("path") == "" {
		rs, err = cm.c.GetRoutes()
	} else {
		rs, _, err = cm.c.Resolve(r.URL.Query().Get("path"))
	}
	if err != nil {
		
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_ROUTES_GETTING", "Can't get routes"})
		return
	}
	routes := []Rule{}
	for _, r := range rs {
		routes = append(routes, Rule{Path: r.Route, Key: r.Key, Kind: r.Name, Storage: r.StorageName})
	}
	_ = json.NewEncoder(w).Encode(&routes)
}

func (cm *configManage) AddRoute(w http.ResponseWriter, r *http.Request) {
	var ruleOptions config.Rule
	err := json.NewDecoder(r.Body).Decode(&ruleOptions)
	if err != nil {
		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %e", err)})
		return
	}
	err = cm.c.AddRule(ruleOptions)
	if err != nil {
		writeError(w, &errorResponse{"WRONG_RULE", "Unnown rule type"})
		return
	}
	w.WriteHeader(http.StatusCreated)
}
