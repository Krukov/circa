package manage

import (
	"circa/config"
	"circa/handler"
	"circa/message"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type runnerHandler struct {
	runner *handler.Runner
}

func newRunnerHandler(r *handler.Runner) *runnerHandler {
	return &runnerHandler{runner: r}
}

type HandlerItem struct {
	ID      uuid.UUID `json:"id"`
	Path    string    `json:"path"`
	Key     string    `json:"key"`
	Rule    string    `json:"type"`
	Storage string    `json:"storage"`
	Methods []string  `json:"methods"`
}

type ProxyTarget struct {
	Target  string   `json:"target"`
	Timeout Duration `json:"timeout"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var handlerPath = regexp.MustCompile("^/api/handler/(.+)$")

func (h *runnerHandler) Handlers(w http.ResponseWriter, r *http.Request) {
	m := handlerPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		h.GetAllHandlers(w, r)
		return
	}
	id, err := uuid.Parse(m[1])
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Handler not found"})
	}
	if r.Method == http.MethodGet {
		h.GetHandler(w, r, id)
		return
	}
	if r.Method == http.MethodDelete {
		h.DelHandler(w, r, id)
		return
	}
	writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
}

func (h *runnerHandler) GetHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	handler := h.runner.GetHandlerByID(id)
	if handler == nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Handler not found"})
		return
	}
	_h := &HandlerItem{
		ID:      handler.ID,
		Path:    handler.Path,
		Key:     handler.Key,
		Rule:    handler.Rule,
		Storage: handler.Storage,
		Methods: handler.Methods,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(_h)
}

func (h *runnerHandler) DelHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	err := h.runner.DelHandlerByID(id)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&errorResponse{"NOT_FOUND", "Handler not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "OK"}`))
}

func (h *runnerHandler) GetAllHandlers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	handlerItems := []*HandlerItem{}
	for _, handler := range h.runner.GetHandlers() {
		handlerItems = append(handlerItems, &HandlerItem{
			ID:      handler.ID,
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

func (h *runnerHandler) GetHandlersByRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.AddRule(w, r)
		return
	}
	if r.Method != http.MethodGet {
		writeError(w, &errorResponse{"METHOD_NOT_ALLOWED", "Method not allowed"})
		return
	}
	handlerItems := []*HandlerItem{}
	method := "get"
	if r.URL.Query().Get("method") != "" {
		method = r.URL.Query().Get("method")
	}
	rule := ""
	if r.URL.Query().Get("type") != "" {
		rule = r.URL.Query().Get("type")
	}
	storage := ""
	if r.URL.Query().Get("storage") != "" {
		storage = r.URL.Query().Get("storage")
	}
	handlers, _ := h.runner.GetHandlersFor(r.URL.Query().Get("path"), method)
	for _, handler := range handlers {
		if rule != "" && string(handler.Rule) != rule {
			continue
		}
		if storage != "" && handler.Storage != storage {
			continue
		}
		handlerItems = append(handlerItems, &HandlerItem{
			ID:      handler.ID,
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
	nh := h.runner.AddHandlers(ruleOptions.Path, handler.NewHandler(rule, storage, ruleOptions.Key, defRequest, ruleOptions.Methods))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var handlerItem *HandlerItem
	for _, handler := range nh {
		handlerItem = &HandlerItem{
			ID:      handler.ID,
			Path:    handler.Path,
			Key:     handler.Key,
			Rule:    handler.Rule,
			Storage: handler.Storage,
			Methods: handler.Methods,
		}
	}
	_ = json.NewEncoder(w).Encode(&handlerItem)
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
