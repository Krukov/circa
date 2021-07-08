package handler

import (
	"circa/message"
	"circa/storages"
	"reflect"
	"testing"
)

func fakeCall(r *message.Request) (*message.Response, error) {
	return &message.Response{Status: 100}, nil
}

type MyRule struct{}

func (r *MyRule) String() string {
	return "my"
}

func (r *MyRule) Process(request *message.Request, key string, storage storages.Storage, call message.Requester) (*message.Response, bool, error) {
	return &message.Response{Status: 200, CachedKey: "test", Headers: map[string]string{}}, false, nil
}

func TestRunner_Handle(t *testing.T) {

	r := NewRunner(fakeCall)
	h := NewHandler(&MyRule{}, &storages.Memory{}, "key", &message.Request{}, []string{"get"})
	r.AddHandlers("/", h)

	tests := []struct {
		name     string
		request  *message.Request
		wantResp *message.Response
		wantErr  bool
	}{
		{"simple", &message.Request{Path: "/", Method: "get"}, &message.Response{Status: 200, Headers: map[string]string{"X-Circa-Cache-Key": "test"}, CachedKey: "test"}, false},
		{"404", &message.Request{Path: "/404", Method: "get"}, &message.Response{Status: 100}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.Handle(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v, resp %v", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Handle() got = %v, want %v", got, tt.wantResp)
			}
		})
	}
}
