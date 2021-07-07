package handler

import (
	"circa/message"
	"circa/rules"
	"circa/storages"
	"strconv"
	"strings"
	"time"
)

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

func (h *handler) ToCall (call message.Requester)  message.Requester {
	return func(request *message.Request) (*message.Response, error) {
		return h.Run(request, call)
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


func (r *Runner) AddHandlers (route string, handlers ...*handler) {
	r.handlers[ruleName(route)] = append(r.handlers[ruleName(route)], handlers...)
	r.router.addRule(route, ruleName(route))
}

func (r *Runner) SetProxy (target string, timeout time.Duration) {
	defRequest := &message.Request{Timeout: timeout, Host: target}
	h := &handler{rule: &rules.ProxyRule{},  defaultRequest: defRequest, methods: ALL_METHODS}
	r.AddHandlers("*", h)
}

func (r *Runner) Handle (request *message.Request) (resp *message.Response, err error) {
	ruleNames, params, err := r.router.resolve(request.Path)
	if err != nil {
		if err == NotFound {
			request.Logger.Debug().Msg("Route for request not found. Forward request")
		} else {
			return nil, err
		}
	}
	request.Logger = request.Logger.With().Str("route", string(ruleNames[0])).Bool("cached", false).Logger()

	makeRequest := r.makeRequest
	for _, rule := range ruleNames {

		handlers_, ok := r.handlers[rule]
		if !ok {
			request.Logger.Warn().Msg("Rule found but no handlers with this rule name")
			return nil, NotFound
		}
		request.Params = params
		for _, handler_ := range handlers_ {
			makeRequest = handler_.ToCall(makeRequest)
		}
	}
	resp, err = makeRequest(request)
	if err != nil {
		return nil, err
	}
	request.Logger = request.Logger.With().Str("status", strconv.Itoa(resp.Status)).Logger()
	if resp.CachedKey != "" {
		request.Logger = request.Logger.With().Bool("cached", true).Str("cacheKey", resp.CachedKey).Logger()
		resp.Headers["X-Circa-Cache-Key"] = resp.CachedKey
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