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
	rule        rules.Rule
	storage     storages.Storage

	defaultRequest *message.Request
	methods        map[string]bool
}

func NewHandler(rule rules.Rule, storage storages.Storage, keyTemplate string, defaultRequest *message.Request, methods []string,
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
		rule:           rule,
		storage:        storage,
		keyTemplate:    keyTemplate,
		defaultRequest: defaultRequest,
		methods:        methodsMap,
	}
}

func (h *handler) ToCall(call message.Requester, route string) message.Requester {
	return func(request *message.Request) (*message.Response, error) {
		resp, hit, err := h.Run(request, call)
		status := "set"
		if err != nil {
			status = "error"
		} else if hit {
			status = "get"
		}
		routeHandlerCount.WithLabelValues(h.rule.String(), route, h.keyTemplate, status).Inc()
		return resp, err
	}
}

func (h *handler) Run(request *message.Request, call message.Requester) (*message.Response, bool, error) {
	if _, ok := h.methods[strings.ToLower(request.Method)]; !ok {
		resp, err := call(request)
		return resp, false, err
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
	request.Params["request_path"] = request.Path
	return formatTemplate(h.keyTemplate, request.Params)
}

type Runner struct {
	handlers    map[ruleName][]*handler
	router      *node
	makeRequest message.Requester
	target      string
	timeout     time.Duration
}

func NewRunner(makeRequest message.Requester) *Runner {
	return &Runner{handlers: map[ruleName][]*handler{}, router: newTrie(), makeRequest: makeRequest}
}

func (r *Runner) AddHandlers(route string, handlers ...*handler) {
	r.handlers[ruleName(route)] = append(r.handlers[ruleName(route)], handlers...)
	r.router.addRule(route, ruleName(route))
	for _, h := range handlers {
		handlersGauge.WithLabelValues(h.rule.String(), route).Inc()
	}
}

func (r *Runner) SetProxy(target string, timeout time.Duration) {
	r.target = target
	r.timeout = timeout
}

func (r *Runner) Handle(request *message.Request) (resp *message.Response, err error) {
	request.Host = r.target
	request.Timeout = r.timeout
	ruleNames, params, err := r.router.resolve(request.Path)
	if err != nil {
		if err == NotFound {
			request.Logger.Debug().Msg("Route for request not found. Forward request")
		} else {
			return nil, err
		}
	}

	makeRequest := r.makeRequest

	for _, rule := range ruleNames {
		request.Logger = request.Logger.With().Str("route", string(rule)).Logger()

		handlers_, ok := r.handlers[rule]
		if !ok {
			request.Logger.Warn().Msg("Rule found but no handlers with this rule name")
			return nil, NotFound
		}
		request.Params = params
		for _, handler_ := range handlers_ {
			makeRequest = handler_.ToCall(makeRequest, string(rule))
		}
	}
	resp, err = makeRequest(request)
	if err != nil {
		return nil, err
	}
	request.Logger = request.Logger.With().Str("status", strconv.Itoa(resp.Status)).Logger()
	if resp.CachedKey != "" {
		request.Logger = request.Logger.With().Str("cacheKey", resp.CachedKey).Logger()
		resp.Headers["X-Circa-Cache-Key"] = resp.CachedKey
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
