package handler

import (
	"strings"

	"github.com/google/uuid"

	"circa/message"
	"circa/rules"
	"circa/storages"
)

var ALL_METHODS = map[string]bool{"get": true, "post": true, "head": true, "put": true, "patch": true, "options": true}

type handler struct {
	id          uuid.UUID
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
		id:             uuid.New(),
		rule:           rule,
		storage:        storage,
		keyTemplate:    keyTemplate,
		defaultRequest: defaultRequest,
		methods:        methodsMap,
	}
}

func (h *handler) ToCall(call message.Requester, route string) message.Requester {
	return func(request *message.Request) (resp *message.Response, err error) {
		var hit bool
		var status string
		if request.Skip {
			resp, err = call(request)
			status = "skip"
		} else {
			resp, hit, err = h.Run(request, call)
			status = "pass"
			if err != nil {
				status = "error"
			} else if hit {
				status = "hit"
			}
		}
		routeHandlerCount.WithLabelValues(h.rule.String(), route, h.keyTemplate, status).Inc()
		return resp, err
	}
}

func (h *handler) Run(request *message.Request, call message.Requester) (*message.Response, bool, error) {
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
	params := map[string]string{}
	for k, v := range request.Params {
		params[k] = v
	}
	for hk, hv := range request.Headers {
		params["H:"+strings.ToLower(hk)] = hv
	}
	params["R:path"] = request.Path
	params["R:method"] = request.Method
	params["R:body"] = string(request.Body)
	return formatTemplate(h.keyTemplate, params)
}

func (h *handler) getMethods() []string {
	methods := []string{}
	for method := range h.methods {
		methods = append(methods, method)
	}
	return methods
}
