package manage

import (
	"circa/config"
	"encoding/json"
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

type Route struct {
	Path    string   `json:"path"`
	Key     string   `json:"key"`
	Rule    string   `json:"type"`
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
		} else {
			writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		}
		return
	}
	name := m[1]
	if r.Method == http.MethodGet {
		cm.GetStorage(w, r, name)
		return
	}
	if r.Method == http.MethodDelete {
		// h.DelStorage(w, r, id)
		return
	}
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

var routePath = regexp.MustCompile("^/api/route/(.+)$")

func (cm *configManage) Routes(w http.ResponseWriter, r *http.Request) {
	m := routePath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		if r.Method == http.MethodGet {
			cm.GetRoutes(w, r)
		} else {
			writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		}
		return
	}
	// name := m[1]
	if r.Method == http.MethodGet {
		// cm.GetRoute(w, r, name)
		return
	}
	if r.Method == http.MethodDelete {
		// h.DelStorage(w, r, id)
		return
	}
	writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
}

func (cm *configManage) GetRoutes(w http.ResponseWriter, r *http.Request) {
	rs, err := cm.c.GetRoutes()
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(&errorResponse{"ERROR_ROUTES_GETTIONG", "Can't get routes"})
	}
	routes := []Route{}
	for _, r := range rs {
		routes = append(routes, Route{Path: r.Route, Key: r.Key, Rule: r.Name, Storage: r.Storage.String()})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&routes)
}

// func (h *runnerHandler) GetHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
// 	handler := h.runner.GetHandlerByID(id)
// 	if handler == nil {
// 		w.Header().Set("content-type", "application/json")
// 		w.WriteHeader(http.StatusNotFound)
// 		_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Handler not found"})
// 		return
// 	}
// 	_h := &HandlerItem{
// 		ID:      handler.ID,
// 		Path:    handler.Path,
// 		Key:     handler.Key,
// 		Rule:    handler.Rule,
// 		Storage: handler.Storage,
// 		Methods: handler.Methods,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(_h)
// }

// func (h *runnerHandler) DelHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
// 	err := h.runner.DelHandlerByID(id)
// 	if err != nil {
// 		w.Header().Set("content-type", "application/json")
// 		w.WriteHeader(http.StatusNotFound)
// 		_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Handler not found"})
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"status": "OK"}`))
// }

// 	handlerItems := []*HandlerItem{}
// 	for _, handler := range h.runner.GetHandlers() {
// 		handlerItems = append(handlerItems, &HandlerItem{
// 			ID:      handler.ID,
// 			Path:    handler.Path,
// 			Key:     handler.Key,
// 			Rule:    handler.Rule,
// 			Storage: handler.Storage,
// 			Methods: handler.Methods,
// 		})
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(&handlerItems)
// }

// func (h *runnerHandler) GetHandlersByRoute(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		h.AddRule(w, r)
// 		return
// 	}
// 	if r.Method != http.MethodGet {
// 		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
// 		return
// 	}
// 	handlerItems := []*HandlerItem{}
// 	method := "get"
// 	if r.URL.Query().Get("method") != "" {
// 		method = r.URL.Query().Get("method")
// 	}
// 	rule := ""
// 	if r.URL.Query().Get("type") != "" {
// 		rule = r.URL.Query().Get("type")
// 	}
// 	storage := ""
// 	if r.URL.Query().Get("storage") != "" {
// 		storage = r.URL.Query().Get("storage")
// 	}
// 	handlers, _ := h.runner.GetHandlersFor(r.URL.Query().Get("path"), method)
// 	for _, handler := range handlers {
// 		if rule != "" && string(handler.Rule) != rule {
// 			continue
// 		}
// 		if storage != "" && handler.Storage != storage {
// 			continue
// 		}
// 		handlerItems = append(handlerItems, &HandlerItem{
// 			ID:      handler.ID,
// 			Path:    handler.Path,
// 			Key:     handler.Key,
// 			Rule:    handler.Rule,
// 			Storage: handler.Storage,
// 			Methods: handler.Methods,
// 		})
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(&handlerItems)
// }

// func (h *runnerHandler) AddRule(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
// 		return
// 	}
// 	var ruleOptions config.Rule
// 	err := json.NewDecoder(r.Body).Decode(&ruleOptions)
// 	if err != nil {
// 		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %e", err)})
// 		return
// 	}
// 	rule, err := config.GetRuleFromOptions(ruleOptions)
// 	if err != nil {
// 		writeError(w, &errorResponse{"WRONG_RULE", "Unnown rule type"})
// 		return
// 	}
// 	storage, err := h.runner.GetStorage(ruleOptions.Storage)
// 	if err != nil {
// 		writeError(w, &errorResponse{"WRONG_STORAGE", "Unnown storage name"})
// 		return
// 	}
// 	defRequest := &message.Request{Host: ruleOptions.Target, Timeout: timeFromString(ruleOptions.Timeout)}
// 	nh := h.runner.AddHandlers(ruleOptions.Path, handler.NewHandler(rule, storage, ruleOptions.Key, defRequest, ruleOptions.Methods))

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	var handlerItem *HandlerItem
// 	for _, handler := range nh {
// 		handlerItem = &HandlerItem{
// 			ID:      handler.ID,
// 			Path:    handler.Path,
// 			Key:     handler.Key,
// 			Rule:    handler.Rule,
// 			Storage: handler.Storage,
// 			Methods: handler.Methods,
// 		}
// 	}
// 	_ = json.NewEncoder(w).Encode(&handlerItem)
// }

// func (h *runnerHandler) Target(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		target, timeout := h.runner.GetProxyOptions()
// 		w.Header().Set("Content-Type", "application/json")
// 		_ = json.NewEncoder(w).Encode(&ProxyTarget{Target: target, Timeout: Duration(timeout)})
// 		return
// 	}
// 	if r.Method != http.MethodPost {
// 		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
// 		return
// 	}
// 	var targetOptions ProxyTarget
// 	err := json.NewDecoder(r.Body).Decode(&targetOptions)
// 	if err != nil {
// 		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %v", err)})
// 		return
// 	}
// 	h.runner.SetProxy(targetOptions.Target, time.Duration(targetOptions.Timeout))
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"status": "OK"}`))
// }

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

// func (r *Runner) AddHandlers(route string, handlers ...*handler) []*HandlerInfo {
// 	r.lock.Lock()
// 	for _, h := range handlers {
// 		if _, ok := r.handlers[ruleName(route)]; !ok {
// 			r.handlers[ruleName(route)] = map[uuid.UUID]*handler{}
// 		}
// 		r.handlers[ruleName(route)][h.id] = h
// 	}
// 	r.router.addRule(route, ruleName(route))
// 	r.lock.Unlock()
// 	handlerItems := []*HandlerInfo{}
// 	for _, h := range handlers {
// 		handlersGauge.WithLabelValues(h.rule.String(), route).Inc()
// 		handlerItems = append(handlerItems, &HandlerInfo{
// 			ID:      h.id,
// 			Path:    route,
// 			Key:     h.keyTemplate,
// 			Rule:    h.rule.String(),
// 			Storage: h.storage.String(),
// 			Methods: h.getMethods(),
// 		})
// 	}
// 	return handlerItems
// }

// type HandlerInfo struct {
// 	ID      uuid.UUID
// 	Path    string
// 	Key     string
// 	Rule    string
// 	Storage string
// 	Methods []string
// }

// func (r *Runner) GetHandlers() []*HandlerInfo {
// 	handlerItems := []*HandlerInfo{}
// 	r.lock.RLock()
// 	defer r.lock.RUnlock()
// 	for path, handlers := range r.handlers {
// 		for _, h := range handlers {
// 			handlerItems = append(handlerItems, &HandlerInfo{
// 				ID:      h.id,
// 				Path:    string(path),
// 				Key:     h.keyTemplate,
// 				Rule:    h.rule.String(),
// 				Storage: h.storage.String(),
// 				Methods: h.getMethods(),
// 			})
// 		}
// 	}
// 	return handlerItems
// }

// func (r *Runner) GetHandlersFor(path, method string) ([]*HandlerInfo, map[string]string) {
// 	handlerItems := []*HandlerInfo{}
// 	r.lock.RLock()
// 	ruleNames, params, err := r.config.Resolve(path)
// 	r.lock.RUnlock()
// 	if err != nil {
// 		return handlerItems, params
// 	}
// 	req := message.Request{Method: method, Path: path, Params: params}
// 	for _, rule := range ruleNames {
// 		r.lock.RLock()
// 		handlers_, ok := r.handlers[rule]
// 		r.lock.RUnlock()
// 		req.Route = string(rule)
// 		if !ok {
// 			continue
// 		}
// 		for _, h := range handlers_ {
// 			if _, ok := h.methods[strings.ToLower(method)]; ok {
// 				handlerItems = append(handlerItems, &HandlerInfo{
// 					ID:      h.id,
// 					Path:    string(rule),
// 					Key:     h.makeKey(&req),
// 					Rule:    h.rule.String(),
// 					Storage: h.storage.String(),
// 					Methods: h.getMethods(),
// 				})
// 			}
// 		}
// 	}
// 	return handlerItems, params
// }

// func (r *Runner) GetHandlerByID(id uuid.UUID) *HandlerInfo {
// 	var ok bool
// 	var h *handler
// 	r.lock.RLock()
// 	defer r.lock.RUnlock()
// 	for rule, handlers := range r.handlers {
// 		h, ok = handlers[id]
// 		if ok {
// 			return &HandlerInfo{
// 				ID:      h.id,
// 				Path:    string(rule),
// 				Key:     h.keyTemplate,
// 				Rule:    h.rule.String(),
// 				Storage: h.storage.String(),
// 				Methods: h.getMethods(),
// 			}
// 		}
// 	}
// 	return nil
// }
