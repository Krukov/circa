package rules

import (
	"circa/key_template"
	"circa/message"
	"circa/storages"
	"strings"
	"sync"

	"github.com/valyala/fastjson"
)

type GlueRule struct {
	Calls map[string]string
}

type glueResp struct {
	name string
	resp *message.Response
	err  error
}

func (r *GlueRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	var wg sync.WaitGroup
	l := &sync.Mutex{}
	responses := []*glueResp{}
	for name, path := range r.Calls {
		wg.Add(1)
		go func(req *message.Request, name string) {
			defer wg.Done()
			resp, _, err := simpleCall(req, call)
			l.Lock()
			defer l.Unlock()
			responses = append(responses, &glueResp{name: name, resp: resp, err: err})
		}(copyRequest(request, path), name)
	}
	wg.Wait()
	resp, err := buildResponse(responses)
	return resp, false, err
}

func copyRequest(req *message.Request, path string) *message.Request {
	path = key_template.FormatTemplate(path, req.Params)
	fp := strings.Replace(req.FullPath, req.Path, path, 1)
	headers := map[string]string{}
	for n, h := range req.Headers {
	    headers[n] = h
	}
	r := message.Request{
		FullPath: fp,
		Path:     path,
		Method:   req.Method,
		QueryStr: req.QueryStr,
		Host:     req.Host,
		Body:     req.Body,
		Route:    req.Route,
		Params:   req.Params,
		Query:    req.Query,
		Headers:  headers,
		Timeout:  req.Timeout,
		Skip:     req.Skip,
		Logger:   req.Logger.With().Logger(),
	}
	return &r
}

var parserPool = fastjson.ParserPool{}

func buildResponse(responses []*glueResp) (*message.Response, error) {
	var v *fastjson.Value
	var o *fastjson.Object
	var err error
	var Res *fastjson.Value
	mainParser := parserPool.Get()
	defer parserPool.Put(mainParser)
	copyParser := parserPool.Get()
	defer parserPool.Put(copyParser)

	for _, r := range responses {
		if r.err != nil {
			// todo skip by settings
			return nil, r.err
		}
		if Res == nil {
			v, err = mainParser.ParseBytes(r.resp.Body)
			if err != nil {
				return nil, err
			}
			Res = v
			continue
		}
		v, err = copyParser.ParseBytes(r.resp.Body)
		if err != nil {
			return nil, err
		}
		o = v.GetObject()
		o.Visit(func(k []byte, vv *fastjson.Value) {
			if string(k) != "id" {
				Res.Set(string(k), vv)
			}
		})
	}

	return message.NewResponse(200, Res.MarshalTo(nil), map[string]string{}), nil
}
