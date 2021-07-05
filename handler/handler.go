package handler

import (
	"circa/message"
	"circa/rules"
	"circa/storages"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_HADDLER_NAME = "_default"
var ALL_METHODS = map[string]bool{"get": true, "post": true, "head": true, "put": true, "patch": true, "options": true}

type handler struct {

	keyTemplate string
	rule rules.Rule
	storage storages.Storage

	defaultRequest *message.Request
	methods map[string]bool


}

func NewHandler (rule rules.Rule, storage storages.Storage, keyTemplate string, defaultRequest *message.Request, methods []string,
) *handler {
	methodsMap := map[string]bool{}
	if len(methods) > 0 {
		for _, method := range methods {
			methodsMap[strings.ToLower(method)] = true
		}
	} else {
		methodsMap["get"] = true
	}
	return &handler{
		rule: rule,
		storage: storage,
		keyTemplate: keyTemplate,
		defaultRequest: defaultRequest,
		methods: methodsMap,
	}
}

func (h *handler) Run (request *message.Request, call message.Requester) (*message.Response, error) {
	if _, ok := h.methods[strings.ToLower(request.Method)]; !ok {
		return call(request)
	}
	request = mergeRequests(request, h.defaultRequest)
	key := h.makeKey(request)
	logger := request.Logger.With().
		Stringer("storage", h.storage).
		Stringer("rule", h.rule).
		Str("key", key).Logger()
	request.Logger = logger
	logger.Debug().Msg("Process rule")
	return h.rule.Process(request, key, h.storage, call)
}

func (h *handler) makeKey(request *message.Request) string {
	return formatTemplate(h.keyTemplate, request.Params)
}

type Runner struct {
	handlers    map[ruleName][]*handler
	router      *node
	makeRequest message.Requester
}

func NewRunner(makeRequest message.Requester) *Runner {
	return &Runner{map[ruleName][]*handler{}, newTrie(), makeRequest}
}


func (r *Runner) AddHandler (route string, handler *handler) {
	r.handlers[ruleName(route)] = append(r.handlers[ruleName(route)], handler)
	r.router.addRule(route, ruleName(route))
}

func (r *Runner) SetProxy (target string, timeout time.Duration) {
	defRequest := &message.Request{Timeout: timeout, Host: target}
	h := &handler{rule: &rules.ProxyRule{},  defaultRequest: defRequest, methods: ALL_METHODS}
	r.handlers[DEFAULT_HADDLER_NAME] = []*handler{h}
}

func (r *Runner) Handle (request *message.Request) (resp *message.Response, err error) {
	ruleName_, params, err := r.router.getRoute(request.Path)
	if err != nil {
		if err == NotFound {
			request.Logger.Debug().Msg("Route for request not found. Forward request")
			ruleName_ = DEFAULT_HADDLER_NAME
		} else {
			return nil, err
		}
	}
	request.Logger = request.Logger.With().Str("route", string(ruleName_)).Bool("cached", false).Logger()
	handlers_, ok := r.handlers[ruleName_]
	if !ok {
		request.Logger.Warn().Msg("Rule found but no handlers with this rule name")
		return nil, NotFound
	}
	request.Params = params
	for _, handler_ := range handlers_ {
		resp, err = handler_.Run(request, r.makeRequest)
		if err != nil {
			return nil, err
		}
		request.Logger = request.Logger.With().Str("status", strconv.Itoa(resp.Status)).Logger()
		if resp.CachedKey != "" {
			request.Logger = request.Logger.With().Bool("cached", true).Str("cacheKey", resp.CachedKey).Logger()
			resp.Headers["Circa-Cache-Key"] = resp.CachedKey
			return
		}
	}
	return
}


func mergeRequests(req1, req2 *message.Request) *message.Request {
	if req1.Timeout.Seconds() == 0.0 {
		req1.Timeout = req2.Timeout
	}
	if req1.Host  == "" {
		req1.Host = req2.Host
	}
	return req1
}