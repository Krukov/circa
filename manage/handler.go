package manage

import (
	"circa/config"
	"circa/rules"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type configManage struct {
	c *config.Config
}

func newConfigManage(c *config.Config) *configManage {
	return &configManage{c: c}
}

type Rule struct {
	Path    string   `json:"path"`
	Key     string   `json:"key"`
	Kind    string   `json:"kind"`
	Storage string   `json:"storage"`
	Methods []string `json:"methods"`
}

type Storage struct {
	Name    string `json:"name"`
	Setting string `json:"setting"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var storagePath = regexp.MustCompile("^/api/storage/(.+)$")

func (cm *configManage) Storages(w http.ResponseWriter, r *http.Request) {
	m := storagePath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		if r.Method == http.MethodGet {
			cm.GetStorages(w, r)
			return
		}
		if r.Method == http.MethodPost {
			cm.CreateStorage(w, r)
			return
		}
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	name := m[1]
	if r.Method == http.MethodGet {
		cm.GetStorage(w, r, name)
		return
	}
	if r.Method == http.MethodPatch {
		cm.ChangeStorage(w, r, name)
		return
	}
	// if r.Method == http.MethodDelete {
	// 	h.DelStorage(w, r, id)
	// 	return
	// }
	writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
}

func (cm *configManage) GetStorages(w http.ResponseWriter, r *http.Request) {
	jsonStorage := []Storage{}
	storages, _ := cm.c.GetStorages()
	for n, s := range storages {
		jsonStorage = append(jsonStorage, Storage{Name: n, Setting: s})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&jsonStorage)
}

func (cm *configManage) CreateStorage(w http.ResponseWriter, r *http.Request) {
	var storage Storage
	err := json.NewDecoder(r.Body).Decode(&storage)
	if err != nil {
		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %e", err)})
		return
	}
	// TODO: Validation settings
	err = cm.c.AddStorage(storage.Name, storage.Setting)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_STORAGE_ADDING", "Can't add storage"})
		return
	}
	if r.URL.Query().Get("sync") != "" {
		cm.c.Sync()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&storage)
}

func (cm *configManage) GetStorage(w http.ResponseWriter, r *http.Request, name string) {
	storages, _ := cm.c.GetStorages()
	for n, s := range storages {
		if n == name {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(&Storage{Name: name, Setting: s})
			return
		}
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Storage not found"})
}

func (cm *configManage) ChangeStorage(w http.ResponseWriter, r *http.Request, name string) {
	var storage Storage
	err := json.NewDecoder(r.Body).Decode(&storage)
	if err != nil {
		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %e", err)})
		return
	}
	// TODO: Validation settings
	err = cm.c.AddStorage(name, storage.Setting)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_STORAGE_ADDING", "Can't add storage"})
		return
	}
	if r.URL.Query().Get("sync") != "" {
		cm.c.Sync()
	}
	cm.GetStorage(w, r, name)
}

var rulePath = regexp.MustCompile("^/api/rules/(.+)$")

func (cm *configManage) Rules(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_ROUTES_GETTING", "Can't get routes"})
		return
	}
	routes := []Rule{}
	for _, r := range rs {
		routes = append(routes, Rule{Path: r.Route, Key: r.Key, Kind: r.Name, Storage: r.StorageName})
	}
	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func timeFromString(in string) time.Duration {
	res, err := time.ParseDuration(in)
	if in == "" || err != nil {
		return time.Second
	}
	return res
}

func writeError(w http.ResponseWriter, err *errorResponse) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(err)
}
