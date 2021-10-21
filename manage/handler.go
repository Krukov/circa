package manage

import (
	"circa/config"
	"circa/handler"
	"circa/message"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type runnerHandler struct {
	runner *handler.Runner
}

func newRunnerHandler(r *handler.Runner) *runnerHandler {
	return &runnerHandler{runner: r}
}

type HandlerItem struct {
	Path    string   `json:"path"`
	Key     string   `json:"key"`
	Rule    string   `json:"type"`
	Storage string   `json:"storage"`
	Methods []string `json:"methods"`
}

type ProxyTarget struct {
	Target  string   `json:"target"`
	Timeout Duration `json:"timeout"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *runnerHandler) GetAllHandlers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	handlerItems := []*HandlerItem{}
	for _, handler := range h.runner.GetHandlers() {
		handlerItems = append(handlerItems, &HandlerItem{
			Path:    handler.Path,
			Key:     handler.Key,
			Rule:    handler.Rule,
			Storage: handler.Storage,
			Methods: handler.Methods,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&handlerItems)
}

func (h *runnerHandler) GetHandlers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	handlerItems := []*HandlerItem{}
	method := "get"
	if r.URL.Query().Get("method") != "" {
		method = r.URL.Query().Get("method")
	}
	handlers, _ := h.runner.GetHandlersFor(r.URL.Query().Get("path"), method)
	for _, handler := range handlers {
		handlerItems = append(handlerItems, &HandlerItem{
			Path:    handler.Path,
			Key:     handler.Key,
			Rule:    handler.Rule,
			Storage: handler.Storage,
			Methods: handler.Methods,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&handlerItems)
}

func (h *runnerHandler) AddRule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	var ruleOptions config.Rule
	err := json.NewDecoder(r.Body).Decode(&ruleOptions)
	if err != nil {
		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %e", err)})
		return
	}
	rule, err := config.GetRuleFromOptions(ruleOptions)
	if err != nil {
		writeError(w, &errorResponse{"WRONG_RULE", "Unnown rule type"})
		return
	}
	storage, err := h.runner.GetStorage(ruleOptions.Storage)
	if err != nil {
		writeError(w, &errorResponse{"WRONG_STORAGE", "Unnown storage name"})
		return
	}
	defRequest := &message.Request{Host: ruleOptions.Target, Timeout: timeFromString(ruleOptions.Timeout)}
	h.runner.AddHandlers(ruleOptions.Path, handler.NewHandler(rule, storage, ruleOptions.Key, defRequest, ruleOptions.Methods))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "OK"}`))
}

func (h *runnerHandler) Target(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		target, timeout := h.runner.GetProxyOptions()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(&ProxyTarget{Target: target, Timeout: Duration(timeout)})
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	var targetOptions ProxyTarget
	err := json.NewDecoder(r.Body).Decode(&targetOptions)
	if err != nil {
		writeError(w, &errorResponse{"FORMAT_ERROR", fmt.Sprintf("Can't decode body %v", err)})
		return
	}
	h.runner.SetProxy(targetOptions.Target, time.Duration(targetOptions.Timeout))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "OK"}`))
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
