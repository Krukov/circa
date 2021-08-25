package manage

import (
	"circa/handler"
	"encoding/json"
	"net/http"
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

func (h *runnerHandler) GetAllHandlers(w http.ResponseWriter, r *http.Request) {
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
