package handler

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"circa/message"
	"circa/storages"

	"github.com/google/uuid"
)

type Runner struct {
	handlers    map[ruleName]map[uuid.UUID]*handler
	storages    map[string]storages.Storage
	router      *node
	makeRequest message.Requester
	target      string
	timeout     time.Duration

	sLock *sync.RWMutex
	lock  *sync.RWMutex
}

func NewRunner(makeRequest message.Requester) *Runner {
	return &Runner{
		handlers:    map[ruleName]map[uuid.UUID]*handler{},
		storages:    map[string]storages.Storage{},
		router:      newTrie(),
		makeRequest: makeRequest,
		lock:        &sync.RWMutex{},
		sLock:       &sync.RWMutex{},
	}
}

func (r *Runner) AddStorage(name string, storage storages.Storage) {
	r.sLock.Lock()
	defer r.sLock.Unlock()
	r.storages[name] = storage
}

func (r *Runner) DelStorage(name string) {
	r.sLock.Lock()
	defer r.sLock.Unlock()
	delete(r.storages, name)
}

func (r *Runner) GetStorage(name string) (storages.Storage, error) {
	r.sLock.RLock()
	defer r.sLock.RUnlock()
	s, ok := r.storages[name]
	if !ok {
		return nil, errors.New("no storage")
	}
	return s, nil
}

func (r *Runner) GetStorages() map[string]storages.Storage {
	r.sLock.RLock()
	defer r.sLock.RUnlock()
	sts := map[string]storages.Storage{}
	for n, s := range r.storages {
		sts[n] = s
	}
	return sts
}

func (r *Runner) AddHandlers(route string, handlers ...*handler) []*HandlerInfo {
	r.lock.Lock()
	for _, h := range handlers {
		if _, ok := r.handlers[ruleName(route)]; !ok {
			r.handlers[ruleName(route)] = map[uuid.UUID]*handler{}
		}
		r.handlers[ruleName(route)][h.id] = h
	}
	r.router.addRule(route, ruleName(route))
	r.lock.Unlock()
	handlerItems := []*HandlerInfo{}
	for _, h := range handlers {
		handlersGauge.WithLabelValues(h.rule.String(), route).Inc()
		handlerItems = append(handlerItems, &HandlerInfo{
			ID:      h.id,
			Path:    route,
			Key:     h.keyTemplate,
			Rule:    h.rule.String(),
			Storage: h.storage.String(),
			Methods: h.getMethods(),
		})
	}
	return handlerItems
}

type HandlerInfo struct {
	ID      uuid.UUID
	Path    string
	Key     string
	Rule    string
	Storage string
	Methods []string
}

func (r *Runner) GetHandlers() []*HandlerInfo {
	handlerItems := []*HandlerInfo{}
	r.lock.RLock()
	defer r.lock.RUnlock()
	for path, handlers := range r.handlers {
		for _, h := range handlers {
			handlerItems = append(handlerItems, &HandlerInfo{
				ID:      h.id,
				Path:    string(path),
				Key:     h.keyTemplate,
				Rule:    h.rule.String(),
				Storage: h.storage.String(),
				Methods: h.getMethods(),
			})
		}
	}
	return handlerItems
}

func (r *Runner) GetHandlersFor(path, method string) ([]*HandlerInfo, map[string]string) {
	handlerItems := []*HandlerInfo{}
	r.lock.RLock()
	ruleNames, params, err := r.router.resolve(path)
	r.lock.RUnlock()
	if err != nil {
		return handlerItems, params
	}
	req := message.Request{Method: method, Path: path, Params: params}
	for _, rule := range ruleNames {
		r.lock.RLock()
		handlers_, ok := r.handlers[rule]
		r.lock.RUnlock()
		req.Route = string(rule)
		if !ok {
			continue
		}
		for _, h := range handlers_ {
			if _, ok := h.methods[strings.ToLower(method)]; ok {
				handlerItems = append(handlerItems, &HandlerInfo{
					ID:      h.id,
					Path:    string(rule),
					Key:     h.makeKey(&req),
					Rule:    h.rule.String(),
					Storage: h.storage.String(),
					Methods: h.getMethods(),
				})
			}
		}
	}
	return handlerItems, params
}

func (r *Runner) GetHandlerByID(id uuid.UUID) *HandlerInfo {
	var ok bool
	var h *handler
	r.lock.RLock()
	defer r.lock.RUnlock()
	for rule, handlers := range r.handlers {
		h, ok = handlers[id]
		if ok {
			return &HandlerInfo{
				ID:      h.id,
				Path:    string(rule),
				Key:     h.keyTemplate,
				Rule:    h.rule.String(),
				Storage: h.storage.String(),
				Methods: h.getMethods(),
			}
		}
	}
	return nil
}

func (r *Runner) DelHandlerByID(id uuid.UUID) error {
	var ok bool
	r.lock.RLock()
	defer r.lock.RUnlock()
	for rule, handlers := range r.handlers {
		_, ok = handlers[id]
		if ok {
			delete(r.handlers[rule], id)
			return nil
		}
	}
	return errors.New("not found")
}

func (r *Runner) SetProxy(target string, timeout time.Duration) {
	r.target = target
	if timeout.Microseconds() > 0.0 {
		r.timeout = timeout
	}
}

func (r *Runner) GetProxyOptions() (string, time.Duration) {
	return r.target, r.timeout
}

func (r *Runner) Handle(request *message.Request) (resp *message.Response, err error) {
	request.Route = "-"
	request.Host = r.target
	request.Timeout = r.timeout
	r.lock.RLock()
	ruleNames, params, err := r.router.resolve(request.Path)
	r.lock.RUnlock()
	if err != nil {
		if err == NotFound {
			request.Logger.Debug().Msg("Route for request not found. Forward request")
		} else {
			return nil, err
		}
	}
	request.Params = params
	makeRequest := r.makeRequest

	for _, rule := range ruleNames {
		r.lock.RLock()
		handlers_, ok := r.handlers[rule]
		r.lock.RUnlock()
		if !ok {
			request.Logger.Warn().Msg("Rule found but no handlers with this rule name")
			return nil, NotFound
		}

		request.Route = string(rule)
		for _, handler_ := range handlers_ {
			if _, ok := handler_.methods[strings.ToLower(request.Method)]; ok {
				request.Logger.Debug().Msgf("%s Add handler %v", rule, handler_.rule.String())
				makeRequest = handler_.ToCall(makeRequest, request.Route)
			}
		}
	}
	resp, err = makeRequest(request)
	if err != nil {
		return nil, err
	}
	request.Logger = request.Logger.With().Str("status", strconv.Itoa(resp.Status)).Logger()
	if resp.CachedKey != "" {
		request.Logger = request.Logger.With().Str("cache_key", resp.CachedKey).Logger()
		resp.SetHeader("X-Circa-Cache-Key", resp.CachedKey)
	}
	return
}

func mergeRequests(req1, req2 *message.Request) *message.Request {
	if req1.Timeout.Seconds() == 0.0 {
		req1.Timeout = req2.Timeout
	}
	if req1.Host == "" {
		req1.Host = req2.Host
	}
	return req1
}
