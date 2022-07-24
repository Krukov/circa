package runner

import (
	"strconv"
	"strings"

	"circa/config"
	"circa/key_template"
	"circa/message"
	"circa/rules"
)

type Runner struct {
	config *config.Config
}

func NewRunner(conf *config.Config) *Runner {
	return &Runner{
		config: conf,
	}
}

func (r *Runner) Handle(request *message.Request, makeRequest message.Requester) (resp *message.Response, err error) {
	rules, params, err := r.config.Resolve(request.Path)
	if err != nil {
		return nil, err
	}
	request.Route = "-"
	request.Host, err = r.config.GetTarget()
	if err != nil {
		return nil, err
	}
	request.Timeout, err = r.config.GetTimeout()
	if err != nil {
		return nil, err
	}
	request.Params = params
	for _, rule := range rules {
		request.Route = rule.Route
		if _, ok := rule.Methods[strings.ToLower(request.Method)]; ok {
			makeRequest = r.toCall(makeRequest, rule)
		}
	}
	resp, err = makeRequest(request)
	if err != nil {
		return nil, err
	}
	request.Logger = request.Logger.With().Str("status", strconv.Itoa(resp.Status)).Logger()
	return
}

func (r *Runner) toCall(call message.Requester, rule *rules.Rule) message.Requester {
	return func(request *message.Request) (resp *message.Response, err error) {
		var hit bool
		var status string
		if request.Skip {
			resp, err = call(request)
			status = "skip"
		} else {
			resp, hit, err = r.run(request, call, rule)
			status = "pass"
			if err != nil {
				status = "error"
			} else if hit {
				status = "hit"
			}
		}
		request.Logger.Debug().Msgf("Status: %s", status)
		routeHandlerCount.WithLabelValues(rule.Name, rule.Route, rule.Key, status).Inc()
		return resp, err
	}
}

func (r *Runner) run(request *message.Request, call message.Requester, rule *rules.Rule) (*message.Response, bool, error) {
	key := r.makeKey(request, rule)
	logger := request.Logger.With().
		Stringer("storage", rule.Storage).
		Str("rule", rule.Name).
		Str("key", key).Logger()
	request.Logger = logger
	logger.Debug().Msg("Process rule")
	return rule.Process(request, key, rule.Storage, call)
}

func (r *Runner) makeKey(request *message.Request, rule *rules.Rule) string {
	params := map[string]string{}
	for k, v := range request.Params {
		params[k] = v
	}
	for hk, hv := range request.Headers {
		params["H:"+strings.ToLower(hk)] = hv
	}
	params["R:path"] = request.Path
	params["R:query"] = request.QueryStr
	params["R:full_path"] = request.FullPath
	params["R:method"] = request.Method
	params["R:body"] = string(request.Body)
	return key_template.FormatTemplate(rule.Key, params)
}
